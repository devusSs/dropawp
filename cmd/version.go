package cmd

import (
	"fmt"

	"github.com/devusSs/dropawp/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Check current version & check and apply updates",
	Run: func(_ *cobra.Command, _ []string) {
		build := version.GetBuild()

		if versionPrintJSON {
			fmt.Println(build.JSON())
			return
		}

		if versionPrintGoString {
			fmt.Println(build.String())
			return
		}

		build.PrettyPrint()
	},
}

var (
	versionPrintJSON     bool
	versionPrintGoString bool
)

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVarP(&versionPrintJSON, "json", "j", false, "Print version as JSON")
	versionCmd.Flags().BoolVarP(&versionPrintGoString, "go", "g", false, "Print version as Go string")
}
