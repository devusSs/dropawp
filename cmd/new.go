package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devusSs/dropawp/internal/stdin"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Initializes a new tracking project",
	Run: func(_ *cobra.Command, _ []string) {
		name, err := stdin.ReadFromStdin("Enter the name of the new tracking project (required)")
		cobra.CheckErr(err)

		if name == "" {
			cobra.CheckErr("Project name cannot be empty")
		}

		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configFile := filepath.Join(homeDir, ".dropawp", "configs", name+".json")
		storageDir := filepath.Join(homeDir, ".dropawp", "storages", name)

		fmt.Println("Creating new tracking project...")

		// TODO: actually implement the logic

		fmt.Println("Created project with:")
		fmt.Printf("  - Name:\t\t%s\n", name)
		fmt.Printf("  - Config file:\t%s\n", configFile)
		fmt.Printf("  - Storage directory:\t%s\n", storageDir)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
