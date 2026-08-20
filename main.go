// Command givenergy-cli is a lightweight CLI that reads live telemetry
// (battery state of charge, solar power in, grid import/export, house load)
// from a GivEnergy inverter over the local network, with no cloud account
// and no extra infrastructure — just the dongle's IP address.
package main

import (
	_ "embed"
	"os"
	"runtime/debug"

	"github.com/doooks/givenergy-cli/cmd"
)

// licenseText is embedded here, rather than in the cmd package, because an
// embed path can't reach outside the embedding file's own directory — and
// LICENSE lives at the module root alongside this file.
//
//go:embed LICENSE
var licenseText string

// version is set at build time via -ldflags "-X main.version=...", locking
// it to the git tag being built (see .github/workflows/release.yml). Left
// unset otherwise so resolveVersion can fall back to the Go toolchain's own
// build-info stamping.
var version string

func main() {
	cmd.LicenseText = licenseText
	cmd.Version = resolveVersion()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// resolveVersion picks the best available version string, in order:
//
//  1. The -ldflags-injected value above, which the release workflow sets to
//     the exact git tag being built.
//  2. `go install github.com/.../givenergy-cli@vX.Y.Z` doesn't go through
//     our build (so ldflags never runs), but the Go toolchain still records
//     the requested module version in the binary's own build info — so
//     that's checked next.
//  3. A plain local `go build`/`go run` from within this git checkout gets
//     neither of the above, but the toolchain still stamps the VCS
//     revision it built from by default; that's the last fallback before
//     giving up and reporting "dev".
func resolveVersion() string {
	if version != "" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	var revision string
	var dirty bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			dirty = s.Value == "true"
		}
	}
	if revision == "" {
		return "dev"
	}
	if len(revision) > 7 {
		revision = revision[:7]
	}
	if dirty {
		return "dev+" + revision + "-dirty"
	}
	return "dev+" + revision
}
