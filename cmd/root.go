package cmd

import (
	"os"

	"github.com/devusSs/dropawp/internal/system"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dropawp",
	Short: "A Go CLI tool to track Counter-Strike inventory prices.",
	Long: `dropawp is a Go CLI tool designed to help you track the prices of Counter-Strike skins.

It allows you to create a configuration for a tracking project, query your inventory automatically,
add missing items via a file, query prices for each marketable item and store them in a file based database.

It is designed to be easy to use and efficient. You simple create a configuration via the command line,
add your secrets when asked for them and then you can start tracking your inventory prices.

The secrets will be stored securely in keyring (system based password store for Windows, Linux and macOS).
The Steam API key is not mandatory in case you do not want to check for required Steam services
before querying your ivnentory. The CSFloat API key howver is required to query the prices of your items.

This tool is in no way affiliated with Valve or Steam or CSFloat.

Make sure to use it with caution and at your own risk. Do not use it for malicious purposes or purposes
that would violate the Steam or CSFloat TOS.`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		err := system.CheckSupported()
		cobra.CheckErr(err)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
