package cmd

import (
	"fmt"

	"github.com/devusSs/dropawp/internal/config"
	"github.com/devusSs/dropawp/internal/secret"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes your dropawp tracking project and needed secrets.",
	Run: func(_ *cobra.Command, _ []string) {
		_, err := config.Read()
		if err == nil && !initOverwriteConfig {
			cobra.CheckErr(
				"Configuration already exists. Use --overwrite-config to overwrite existing configuration.",
			)
		}

		var cfg *config.Config

		switch {
		case initUseEnv:
			config.SetEnvFile(initEnvFile)
			cfg, err = config.FromEnv()
			cobra.CheckErr(err)
		case initUseFile:
			config.SetFile(initFile)
			cfg, err = config.FromFile()
			cobra.CheckErr(err)
		default:
			cfg, err = config.FromInput()
			cobra.CheckErr(err)
		}

		if !cfg.SkipSteamServicesCheck || !cfg.SkipSteamUserCheck {
			err = checkOrInsertSecret(secret.SteamAPIKey, initOverwriteSecrets)
			cobra.CheckErr(err)
		}

		err = checkOrInsertSecret(secret.CSFloatAPIKey, initOverwriteSecrets)
		cobra.CheckErr(err)

		err = config.Write(cfg)
		cobra.CheckErr(err)
	},
}

var (
	initOverwriteConfig  bool
	initUseEnv           bool
	initEnvFile          string
	initUseFile          bool
	initFile             string
	initOverwriteSecrets bool
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().
		BoolVar(&initOverwriteConfig, "overwrite-config", false, "Overwrite existing configuration files")
	initCmd.Flags().
		BoolVar(&initUseEnv, "use-env", false, "Use environment variables for configuration")
	initCmd.Flags().
		StringVar(&initEnvFile, "env-file", "", "Path to the environment file if desired")
	initCmd.Flags().
		BoolVar(&initUseFile, "use-file", false, "Use a file for configuration")
	initCmd.Flags().
		StringVar(&initFile, "file", "", "Path to the configuration file if desired")
	initCmd.Flags().
		BoolVar(&initOverwriteSecrets, "overwrite-secrets", false, "Overwrite existing secrets")

	initCmd.MarkFlagsMutuallyExclusive("use-env", "use-file")
	initCmd.MarkFlagsRequiredTogether("use-file", "file")
}

func checkOrInsertSecret(key secret.Key, overwrite bool) error {
	exists, err := secret.Exists(key)
	if err != nil {
		return fmt.Errorf("failed to check secret %s: %w", key, err)
	}

	if !exists || overwrite {
		var value string
		value, err = secret.GetInput(fmt.Sprintf("Enter value for secret %s", key))
		if err != nil {
			return fmt.Errorf("failed to get input for secret %s: %w", key, err)
		}

		err = secret.Save(key, value)
		if err != nil {
			return fmt.Errorf("failed to save secret %s: %w", key, err)
		}
	}

	return nil
}
