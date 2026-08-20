// Command givenergy-monitor is a lightweight CLI that reads live telemetry
// (battery state of charge, solar power in, grid import/export, house load)
// from a GivEnergy inverter over the local network, with no cloud account
// and no extra infrastructure — just the dongle's IP address.
package main

import (
	_ "embed"
	"os"

	"givenergy-cli/cmd"
)

// licenseText is embedded here, rather than in the cmd package, because an
// embed path can't reach outside the embedding file's own directory — and
// LICENSE lives at the module root alongside this file.
//
//go:embed LICENSE
var licenseText string

func main() {
	cmd.LicenseText = licenseText
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
