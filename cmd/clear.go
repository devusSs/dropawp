package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devusSs/dropawp/internal/stdin"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears all data saved by the application",
	Run: func(_ *cobra.Command, _ []string) {
		confirm, err := stdin.ReadFromStdin(
			"Are you sure you want to clear all data? This action cannot be undone (yes/no)",
		)
		cobra.CheckErr(err)

		if confirm != "yes" {
			fmt.Println("Cancelled, no data was cleared.")
			return
		}

		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		baseDir := filepath.Join(homeDir, ".dropawp")

		err = os.RemoveAll(baseDir)
		cobra.CheckErr(err)

		fmt.Println("All data cleared successfully.")
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
