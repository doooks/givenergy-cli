package givenergy

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// errHeartbeat is returned by parseTransparentResponse when the frame read
// off the wire was a heartbeat request rather than our expected register
// response. The dongle sends these unprompted every few minutes on a
// long-lived connection; callers should skip them and keep reading.
var errHeartbeat = errors.New("received heartbeat frame")

type registerResponse struct {
	deviceAddress byte
	baseRegister  uint16
	values        []uint16
	crcOK         bool
}

// readFrame reads exactly one GivEnergy frame from r. It trusts that the
// stream is aligned on a frame boundary, which holds for the short-lived,
// one-request-per-connection pattern this client uses.
func readFrame(r frameReader) ([]byte, error) {
	header := make([]byte, 6)
	if _, err := readFull(r, header); err != nil {
		return nil, fmt.Errorf("reading frame header: %w", err)
	}
	tid := binary.BigEndian.Uint16(header[0:2])
	if tid != headerTID {
		return nil, fmt.Errorf("unexpected transaction id 0x%04x (wrong protocol, or stream out of sync)", tid)
	}
	bodyLen := binary.BigEndian.Uint16(header[4:6])
	body := make([]byte, bodyLen)
	if _, err := readFull(r, body); err != nil {
		return nil, fmt.Errorf("reading frame body (%d bytes): %w", bodyLen, err)
	}
	return append(header, body...), nil
}

// frameReader is the minimal surface readFrame needs; satisfied by net.Conn.
type frameReader interface {
	Read(p []byte) (n int, err error)
}

func readFull(r frameReader, buf []byte) (int, error) {
	total := 0
	for total < len(buf) {
		n, err := r.Read(buf[total:])
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// parseTransparentResponse decodes a GivEnergy "Transparent" read-input-registers
// response frame (as produced by buildReadInputRegistersRequest's counterpart
// on the dongle).
func parseTransparentResponse(frame []byte) (registerResponse, error) {
	if len(frame) < 8 {
		return registerResponse{}, fmt.Errorf("frame too short: %d bytes", len(frame))
	}
	pid := binary.BigEndian.Uint16(frame[2:4])
	if pid != headerPID {
		return registerResponse{}, fmt.Errorf("unexpected protocol id 0x%04x", pid)
	}
	uid := frame[6]
	if uid != 0x00 && uid != 0x01 {
		return registerResponse{}, fmt.Errorf("unexpected unit id 0x%02x", uid)
	}
	mainFunc := frame[7]
	if mainFunc == mainFuncHeartbeat {
		return registerResponse{}, errHeartbeat
	}
	if mainFunc != mainFuncTransparent {
		return registerResponse{}, fmt.Errorf("unexpected main function code 0x%02x", mainFunc)
	}

	const preludeLen = 8 + 10 + 8 // mbap header + data-adapter serial + padding
	if len(frame) < preludeLen+2 {
		return registerResponse{}, fmt.Errorf("frame too short for transparent prelude: %d bytes", len(frame))
	}
	p := preludeLen
	deviceAddress := frame[p]
	p++
	funcCode := frame[p]
	p++
	isError := funcCode&transparentErrorFlag != 0
	funcCode &^= transparentErrorFlag
	if isError {
		return registerResponse{}, fmt.Errorf("device returned an error response (transparent function 0x%02x)", funcCode)
	}
	if funcCode != transparentFuncReadInput {
		return registerResponse{}, fmt.Errorf("unexpected transparent function 0x%02x", funcCode)
	}

	p += 10 // inverter serial number, not needed
	if len(frame) < p+4 {
		return registerResponse{}, fmt.Errorf("frame too short for register header: %d bytes", len(frame))
	}
	baseRegister := binary.BigEndian.Uint16(frame[p : p+2])
	p += 2
	registerCount := binary.BigEndian.Uint16(frame[p : p+2])
	p += 2

	need := p + int(registerCount)*2 + 2
	if len(frame) < need {
		return registerResponse{}, fmt.Errorf("frame too short for %d registers: have %d bytes, need %d", registerCount, len(frame), need)
	}
	values := make([]uint16, registerCount)
	for i := range values {
		values[i] = binary.BigEndian.Uint16(frame[p : p+2])
		p += 2
	}
	wireCRC := frame[p : p+2]

	crcPayload := frame[26 : len(frame)-2]
	raw := crc16Modbus(crcPayload)
	crcOK := wireCRC[0] == byte(raw) && wireCRC[1] == byte(raw>>8)

	return registerResponse{
		deviceAddress: deviceAddress,
		baseRegister:  baseRegister,
		values:        values,
		crcOK:         crcOK,
	}, nil
}
