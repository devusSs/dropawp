package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/devusSs/dropawp/internal/version"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check and apply available updates",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		var err error
		baseDir, err = os.UserHomeDir()
		cobra.CheckErr(err)
	},
}

var checkUpdateCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for available updates",
	Run: func(_ *cobra.Command, _ []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		available, err := version.CheckForUpdate(ctx)
		cobra.CheckErr(err)

		if !available {
			fmt.Println("No updates available.")
			return
		}

		uaPath := filepath.Join(baseDir, ".dropawp", ".update_available")

		f, err := os.Create(uaPath)
		cobra.CheckErr(err)
		defer f.Close()

		_, err = f.WriteString("true")
		cobra.CheckErr(err)

		fmt.Println("An update is available. Run 'dropawp update apply' to install it.")
	},
}

var applyUpdateCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply an update if available",
	Run: func(_ *cobra.Command, _ []string) {
		uaPath := filepath.Join(baseDir, ".dropawp", ".update_available")

		_, err := os.Stat(uaPath)
		if err != nil {
			if !os.IsNotExist(err) {
				cobra.CheckErr(err)
			}

			fmt.Println("Run 'dropawp update check' to check for updates first.")
			return
		}

		fmt.Println("Applying update...")

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		err = version.GetAndApplyUpdate(ctx)
		cobra.CheckErr(err)

		err = os.Remove(uaPath)
		cobra.CheckErr(err)

		fmt.Println("Update applied successfully. Please restart the application.")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.AddCommand(checkUpdateCmd)
	updateCmd.AddCommand(applyUpdateCmd)
}
