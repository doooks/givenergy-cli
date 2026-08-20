package cmd

import "github.com/spf13/cobra"

// Version holds the build version. It's set by main() from a value injected
// at build time via -ldflags rather than hardcoded here — see main.go.
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version number",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println(Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
