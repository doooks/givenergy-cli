package givenergy

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// DefaultDeviceAddress is the address used to query the inverter's core
// telemetry bank (input registers 0-59). Confirmed against a real captured
// request in dewet22/givenergy-modbus (ReadHolding(0x11,0,60)).
const DefaultDeviceAddress = 0x11

// Client polls a GivEnergy dongle for a live telemetry snapshot over its
// proprietary Modbus-TCP dialect (see frame.go/response.go).
type Client struct {
	Host          string
	Port          int
	DeviceAddress byte
	Timeout       time.Duration
}

// NewClient returns a Client with sane defaults; Port/DeviceAddress/Timeout
// can be overridden on the returned value before use.
func NewClient(host string) *Client {
	return &Client{
		Host:          host,
		Port:          8899,
		DeviceAddress: DefaultDeviceAddress,
		Timeout:       5 * time.Second,
	}
}

// ReadSnapshot opens a fresh TCP connection, requests the core telemetry
// register bank, and closes the connection again. A short-lived
// connection-per-poll avoids having to implement the dongle's heartbeat
// protocol, which only matters for long-lived connections.
func (c *Client) ReadSnapshot() (Snapshot, error) {
	addr := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	conn, err := net.DialTimeout("tcp", addr, c.Timeout)
	if err != nil {
		return Snapshot{}, fmt.Errorf("connecting to %s: %w", addr, err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(c.Timeout)); err != nil {
		return Snapshot{}, fmt.Errorf("setting deadline: %w", err)
	}

	req := buildReadInputRegistersRequest(c.DeviceAddress, 0, 60)
	if _, err := conn.Write(req); err != nil {
		return Snapshot{}, fmt.Errorf("sending request: %w", err)
	}

	// The dongle can interleave an unsolicited heartbeat frame; skip past
	// any we see while waiting for our register response.
	for {
		frame, err := readFrame(conn)
		if err != nil {
			return Snapshot{}, fmt.Errorf("reading response from %s: %w", addr, err)
		}
		resp, err := parseTransparentResponse(frame)
		if err == errHeartbeat {
			continue
		}
		if err != nil {
			return Snapshot{}, fmt.Errorf("parsing response from %s: %w", addr, err)
		}
		return snapshotFromRegisters(resp.values), nil
	}
}
