package cmd

import (
	"fmt"

	"github.com/devusSs/dropawp/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the build information of the application and exits.",
	Run: func(cmd *cobra.Command, args []string) {
		b := version.GetBuild()

		if versionPrintJSON {
			fmt.Println(b.JSON())
			return
		}

		b.PrettyPrint()
	},
}

var versionPrintJSON bool

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVarP(&versionPrintJSON, "json", "j", false, "Output version information in JSON format")
}
