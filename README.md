# givenergy-monitor

A tiny CLI that reads live telemetry from a GivEnergy inverter over your
local network — battery state of charge, solar power in, grid
import/export, house load — with no cloud account, no Docker, no
Home Assistant. It talks directly to the inverter's WiFi/ethernet dongle
using GivEnergy's own Modbus-TCP dialect (see [NOTICE.md](NOTICE.md) for
where that protocol knowledge came from). The protocol layer itself
(`internal/givenergy/`) is still pure Go standard library; the only
dependency is [Cobra](https://github.com/spf13/cobra) for argument parsing.

Copyright (c) 2026 Dan Dukeson <dandukeson@gmail.com>. Licensed under the
[MIT License](LICENSE). Contributions welcome — see
[CONTRIBUTING.md](CONTRIBUTING.md).

## Build

```sh
go build -o givenergy-monitor .
```

Cross-compile for another machine (e.g. a Raspberry Pi) with `GOOS`/`GOARCH`:

```sh
GOOS=linux GOARCH=arm64 go build -o givenergy-monitor-arm64 .
```

Pushing a `v*` tag (e.g. `v1.0.0`) triggers
[`.github/workflows/release.yml`](.github/workflows/release.yml), which
vets, tests, cross-builds binaries for linux/darwin/windows on amd64/arm64,
and attaches them (plus a checksums file) to a GitHub release for that tag.

## Usage

Find your dongle's local IP first (your router's DHCP client list, or the
device settings in the GivEnergy app).

```sh
# one-shot snapshot
./givenergy-monitor --host 192.168.1.50

# live view, refreshing every 5s, until Ctrl+C
./givenergy-monitor --host 192.168.1.50 --watch
```

`--host` can also be set via the `GIVENERGY_HOST` environment variable.

Flags:

| Flag | Default | Meaning |
|---|---|---|
| `--host` | (`GIVENERGY_HOST`) | dongle IP address or hostname |
| `--port` | `8899` | dongle Modbus TCP port |
| `--addr` | `17` (0x11) | GivEnergy device address to query (0x11 = inverter) |
| `--watch` | off | keep polling and redraw in place |
| `--interval` | `5s` | poll interval in watch mode |
| `--timeout` | `5s` | per-poll network timeout |

`givenergy-monitor warranty` prints the software's no-warranty disclaimer,
`givenergy-monitor license` prints the full [LICENSE](LICENSE) text, and
`givenergy-monitor version` prints the build version.

Release binaries have their version locked to the git tag they were built
from, via `-ldflags "-X main.version=..."` in the release workflow. A plain
local `go build` reports `dev`; to stamp a specific version yourself:

```sh
go build -ldflags "-X main.version=$(git describe --tags --always --dirty)" -o givenergy-monitor .
```

## A couple of things worth knowing

- **Grid/battery sign convention.** Confirmed against a live system by
  cross-checking the energy balance (solar ≈ load + grid export + battery
  charge power) on a real reading: positive `p_grid_out` = exporting,
  and — contrary to what the register name suggests — positive raw
  `p_battery` is *charging* the battery, not discharging, so
  `metrics.go` negates it before exposing `Snapshot.BatteryWatts` (positive
  = charging, to match its doc comment). If your numbers ever look
  backwards, that's the place to check.
- **Polling frequency.** The upstream protocol notes that address `0x11`
  responses get forwarded to the GivEnergy cloud, and recommend avoiding
  sub-5-minute polling on it for that reason. The default 5s `--watch`
  interval is convenient but more aggressive than that guidance; if you see
  errors or slowdowns, pass a longer `--interval`.

## How it works

Each poll opens a fresh TCP connection, sends one "read input registers
0-59" request, reads the response, and closes the connection — no
persistent connection, so there's no need to handle the dongle's periodic
heartbeat protocol. See `internal/givenergy/` for the frame encoding and
`internal/givenergy/metrics.go` for the specific register offsets used.
