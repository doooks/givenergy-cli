# givenergy-monitor

A tiny CLI that reads live telemetry from a GivEnergy inverter over your
local network — battery state of charge, solar power in, grid
import/export, house load — with no cloud account, no Docker, no
Home Assistant. 

It talks directly to the inverter's WiFi/ethernet dongle.

Copyright (c) 2026 Dan Dukeson <dandukeson@gmail.com>. Licensed under the
[MIT License](LICENSE).

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

## Example Output
```
Battery SOC    84%
Solar in        975 W
Grid              1 W  (import)
House load      633 W
Battery         299 W  (charging)
```

## How it works

Each poll opens a fresh TCP connection, sends one "read input registers
0-59" request, reads the response, and closes the connection — no
persistent connection, so there's no need to handle the dongle's periodic
heartbeat protocol. See `internal/givenergy/` for the frame encoding and
`internal/givenergy/metrics.go` for the specific register offsets used.
