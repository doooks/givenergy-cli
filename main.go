// Command givenergy-monitor is a lightweight CLI that reads live telemetry
// (battery state of charge, solar power in, grid import/export, house load)
// from a GivEnergy inverter over the local network, with no cloud account
// and no extra infrastructure — just the dongle's IP address.
package main

import (
	"os"

	"givenergy-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
