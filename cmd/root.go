// Package cmd defines the givenergy-cli command-line interface.
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"givenergy-cli/internal/display"
	"givenergy-cli/internal/givenergy"
)

var (
	host     string
	port     int
	addr     int
	watch    bool
	interval time.Duration
	timeout  time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "givenergy-cli",
	Short: "Read live telemetry from a GivEnergy inverter over the local network",
	Long: "givenergy-cli reads battery state of charge, solar power in, grid\n" +
		"import/export, and house load from a GivEnergy inverter's dongle over\n" +
		"local Modbus TCP — no cloud account, no extra infrastructure.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().NFlag() == 0 {
			if err := cmd.Help(); err != nil {
				return err
			}
			cmd.Println()
			cmd.Println(warrantyText)
			return nil
		}
		if host == "" {
			return fmt.Errorf("missing --host (or GIVENERGY_HOST)")
		}
		if addr < 0 || addr > 0xFF {
			return fmt.Errorf("--addr must be between 0 and 255")
		}

		client := givenergy.NewClient(host)
		client.Port = port
		client.DeviceAddress = byte(addr)
		client.Timeout = timeout

		if watch {
			return runWatch(client, interval)
		}
		return runOnce(client)
	},
}

// Execute runs the root command; the returned error has already been
// printed by Cobra, so callers just need to decide the process exit code.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	flags := rootCmd.Flags()
	flags.StringVar(&host, "host", os.Getenv("GIVENERGY_HOST"), "dongle IP address or hostname (or set GIVENERGY_HOST)")
	flags.IntVar(&port, "port", 8899, "dongle Modbus TCP port")
	flags.IntVar(&addr, "addr", givenergy.DefaultDeviceAddress, "GivEnergy device address to query (default is the inverter)")
	flags.BoolVarP(&watch, "watch", "w", false, "keep polling and refresh the display in place, until Ctrl+C")
	flags.DurationVar(&interval, "interval", 5*time.Second, "poll interval in watch mode")
	flags.DurationVar(&timeout, "timeout", 5*time.Second, "per-poll network timeout")
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
