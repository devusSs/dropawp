package cmd

import (
	"os"

	"github.com/devusSs/dropawp/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   version.AppName,
	Short: version.AppDescriptionShort,
	Long:  version.AppDescriptionLong,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
