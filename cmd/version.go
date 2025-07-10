package cmd

import (
	"fmt"

	"github.com/devusSs/dropawp/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version / build and exits.",
	Run: func(_ *cobra.Command, _ []string) {
		b := version.GetBuild()

		if versionPrintJSON {
			fmt.Println(b.JSON())
			return
		}

		if versionPrintGoString {
			fmt.Println(b.String())
			return
		}

		fmt.Println(b.Pretty())
	},
}

var (
	versionPrintJSON     bool
	versionPrintGoString bool
)

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().
		BoolVarP(&versionPrintJSON, "json", "j", false, "Print version in JSON format")
	versionCmd.Flags().
		BoolVarP(&versionPrintGoString, "go-string", "g", false, "Print version in Go string format")
}
