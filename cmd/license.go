package cmd

import "github.com/spf13/cobra"

// LicenseText holds the full license text. It's set by main() from the
// embedded LICENSE file rather than duplicated here — see main.go.
var LicenseText string

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Show the software license",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Print(LicenseText)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(licenseCmd)
}
