// Command givenergy-cli is a lightweight CLI that reads live telemetry
// (battery state of charge, solar power in, grid import/export, house load)
// from a GivEnergy inverter over the local network, with no cloud account
// and no extra infrastructure — just the dongle's IP address.
package main

import (
	_ "embed"
	"os"

	"github.com/doooks/givenergy-cli/cmd"
)

// licenseText is embedded here, rather than in the cmd package, because an
// embed path can't reach outside the embedding file's own directory — and
// LICENSE lives at the module root alongside this file.
//
//go:embed LICENSE
var licenseText string

// version is set at build time via -ldflags "-X main.version=...", locking
// it to the git tag being built (see .github/workflows/release.yml). A
// plain `go build`/`go run` without that flag leaves it as "dev".
var version = "dev"

func main() {
	cmd.LicenseText = licenseText
	cmd.Version = version
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
