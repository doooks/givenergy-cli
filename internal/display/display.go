// Package display formats and prints GivEnergy snapshots to the terminal.
package display

import (
	"fmt"
	"io"
	"time"

	"github.com/doooks/givenergy-cli/internal/givenergy"
)

// FormatSnapshot renders a snapshot as a short, aligned block of text.
func FormatSnapshot(s givenergy.Snapshot) string {
	grid, gridLabel := signedLabel(s.GridWatts, "export", "import")
	battery, batteryLabel := signedLabel(s.BatteryWatts, "charging", "discharging")

	return fmt.Sprintf(
		"Battery SOC   %3d%%\n"+
			"Solar in      %5d W\n"+
			"Grid          %5d W  (%s)\n"+
			"House load    %5d W\n"+
			"Battery       %5d W  (%s)\n",
		s.BatterySOC,
		s.SolarWatts,
		grid, gridLabel,
		s.LoadWatts,
		battery, batteryLabel,
	)
}

// signedLabel splits a signed watts value into a magnitude and a direction
// label. Positive is assumed to mean posLabel (e.g. grid export, battery
// charging) based on the register names (p_grid_out, p_battery) — confirm
// this against your real values the first time you run the tool, and flip
// the assumption in metrics.go if it reads backwards for your system.
func signedLabel(w int, posLabel, negLabel string) (int, string) {
	if w < 0 {
		return -w, negLabel
	}
	return w, posLabel
}

// PrintSnapshot writes a single formatted snapshot to w.
func PrintSnapshot(w io.Writer, s givenergy.Snapshot) {
	fmt.Fprint(w, FormatSnapshot(s))
}

// PrintWatchFrame clears the terminal and prints a snapshot plus a status
// line, for use inside a repeated polling loop.
func PrintWatchFrame(w io.Writer, s givenergy.Snapshot, host string, at time.Time) {
	fmt.Fprint(w, "\x1b[H\x1b[2J") // move cursor home + clear screen
	fmt.Fprintf(w, "GivEnergy @ %s — %s\n\n", host, at.Format("15:04:05"))
	fmt.Fprint(w, FormatSnapshot(s))
	fmt.Fprint(w, "\n(Ctrl+C to exit)\n")
}

// PrintWatchError prints a non-fatal error line inside the watch loop
// without clearing the screen, so the last-good reading stays visible.
func PrintWatchError(w io.Writer, err error, at time.Time) {
	fmt.Fprintf(w, "[%s] read failed: %v\n", at.Format("15:04:05"), err)
}
