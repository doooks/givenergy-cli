// Command givenergy-monitor is a lightweight CLI that reads live telemetry
// (battery state of charge, solar power in, grid import/export, house load)
// from a GivEnergy inverter over the local network, with no cloud account
// and no extra infrastructure — just the dongle's IP address.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"givenergy-cli/internal/display"
	"givenergy-cli/internal/givenergy"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "givenergy-monitor:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("givenergy-monitor", flag.ContinueOnError)
	host := fs.String("host", os.Getenv("GIVENERGY_HOST"), "dongle IP address or hostname (or set GIVENERGY_HOST)")
	port := fs.Int("port", 8899, "dongle Modbus TCP port")
	addr := fs.Int("addr", givenergy.DefaultDeviceAddress, "GivEnergy device address to query (default is the inverter)")
	watch := fs.Bool("watch", false, "keep polling and refresh the display in place, until Ctrl+C")
	interval := fs.Duration("interval", 5*time.Second, "poll interval in watch mode")
	timeout := fs.Duration("timeout", 5*time.Second, "per-poll network timeout")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: givenergy-monitor --host <dongle-ip> [flags]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *host == "" {
		fs.Usage()
		return fmt.Errorf("missing --host (or GIVENERGY_HOST)")
	}
	if *addr < 0 || *addr > 0xFF {
		return fmt.Errorf("--addr must be between 0 and 255")
	}

	client := givenergy.NewClient(*host)
	client.Port = *port
	client.DeviceAddress = byte(*addr)
	client.Timeout = *timeout

	if *watch {
		return runWatch(client, *interval)
	}
	return runOnce(client)
}

func runOnce(client *givenergy.Client) error {
	snap, err := client.ReadSnapshot()
	if err != nil {
		return err
	}
	display.PrintSnapshot(os.Stdout, snap)
	return nil
}

func runWatch(client *givenergy.Client, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	poll := func() {
		snap, err := client.ReadSnapshot()
		if err != nil {
			display.PrintWatchError(os.Stdout, err, time.Now())
			return
		}
		display.PrintWatchFrame(os.Stdout, snap, client.Host, time.Now())
	}

	poll()
	for range ticker.C {
		poll()
	}
	return nil
}
