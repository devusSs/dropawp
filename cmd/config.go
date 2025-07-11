package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/devusSs/dropawp/internal/config"
	"github.com/devusSs/dropawp/internal/secret"
	"github.com/spf13/cobra"
)

var cfg *config.Config

var (
	configCmdShowExtended  bool
	configCmdShowSensitive bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show, edit or delete the configuration.",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		var err error
		cfg, err = config.Read()
		cobra.CheckErr(err)
	},
	Run: func(_ *cobra.Command, _ []string) {
		err := printConfigAsTable(configCmdShowExtended)
		cobra.CheckErr(err)

		if configCmdShowSensitive {
			fmt.Println()

			w := tabwriter.NewWriter(os.Stdout, 0, 0, tabwriterPadding, ' ', 0)
			_, err = fmt.Fprintln(w, "Secret Name\tSecret Value")
			cobra.CheckErr(err)

			_, err = fmt.Fprintln(w, "-----------\t------------")
			cobra.CheckErr(err)

			var csFloatAPIKey string
			csFloatAPIKey, err = getSecret(secret.CSFloatAPIKey)
			cobra.CheckErr(err)

			_, err = fmt.Fprintln(w, "CSFloat API Key\t"+csFloatAPIKey)
			cobra.CheckErr(err)

			if !cfg.SkipSteamServicesCheck || !cfg.SkipSteamUserCheck {
				var steamAPIKey string
				steamAPIKey, err = getSecret(secret.SteamAPIKey)
				cobra.CheckErr(err)

				_, err = fmt.Fprintln(w, "Steam API Key\t"+steamAPIKey)
				cobra.CheckErr(err)
			}

			err = w.Flush()
			cobra.CheckErr(err)
		}
	},
}

var configDeleteCmdSecrets bool

var configDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the configuration.",
	Run: func(_ *cobra.Command, _ []string) {
		if configDeleteCmdSecrets {
			err := secret.DeleteAll()
			cobra.CheckErr(err)

			fmt.Println("All secrets have been deleted.")
		}

		err := config.Delete(cfg)
		cobra.CheckErr(err)

		fmt.Println("Configuration has been deleted successfully.")
	},
}

var (
	configEditProjectName        string
	configEditCooldown           string
	configEditSteamID64          string
	configEditItemsFile          string
	configEditSkipSteamServices  string
	configEditSkipSteamUser      string
	configEditSkipFilterItems    string
	configEditUpdateSecretKeys   []string
	configEditUpdateSecretValues []string
)

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration values using flags.",
	Run: func(_ *cobra.Command, _ []string) {
		updated := false

		if configEditProjectName != "" {
			cfg.ProjectName = configEditProjectName
			updated = true
		}

		if configEditCooldown != "" {
			duration, err := parseExtendedDuration(configEditCooldown)
			cobra.CheckErr(err)

			cfg.CooldownDuration = duration
			updated = true
		}

		if configEditSteamID64 != "" {
			steamID, err := parseUint64(configEditSteamID64)
			cobra.CheckErr(err)

			cfg.SteamID64 = steamID
			updated = true
		}

		if configEditItemsFile != "" {
			cfg.AdditionalItemsFile = configEditItemsFile
			updated = true
		}

		if configEditSkipSteamServices != "" {
			skip, err := parseBool(configEditSkipSteamServices)
			cobra.CheckErr(err)

			cfg.SkipSteamServicesCheck = skip
			updated = true
		}

		if configEditSkipSteamUser != "" {
			skip, err := parseBool(configEditSkipSteamUser)
			cobra.CheckErr(err)

			cfg.SkipSteamUserCheck = skip
			updated = true
		}

		if configEditSkipFilterItems != "" {
			skip, err := parseBool(configEditSkipFilterItems)
			cobra.CheckErr(err)

			cfg.SkipFilterUntradableItems = skip
			updated = true
		}

		if !updated {
			cobra.CheckErr("No changes specified. Use --help to see available flags.")
		}

		cfg.UpdatedAt = time.Now()

		err := config.Write(cfg)
		cobra.CheckErr(err)

		fmt.Println("Configuration updated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().
		BoolVar(&configCmdShowExtended, "extended", false,
			"Show extended configuration values")
	configCmd.Flags().
		BoolVar(&configCmdShowSensitive, "show-sensitive", false, "Show sensitive configuration values")

	configCmd.AddCommand(configDeleteCmd)

	configDeleteCmd.Flags().
		BoolVar(&configDeleteCmdSecrets, "secrets", false,
			"Delete secrets in addition to the configuration file")

	configCmd.AddCommand(configEditCmd)

	configEditCmd.Flags().StringVar(&configEditProjectName, "project-name", "", "Set project name")
	configEditCmd.Flags().
		StringVar(&configEditCooldown, "cooldown", "", "Set cooldown duration (e.g., 30m, 1h, 2d)")
	configEditCmd.Flags().StringVar(&configEditSteamID64, "steam-id", "", "Set Steam ID64")
	configEditCmd.Flags().
		StringVar(&configEditItemsFile, "items-file", "", "Set additional items file path")
	configEditCmd.Flags().
		StringVar(&configEditSkipSteamServices, "skip-steam-services", "", "Skip Steam services check (true/false)")
	configEditCmd.Flags().
		StringVar(&configEditSkipSteamUser, "skip-user-check", "", "Skip Steam user check (true/false)")
	configEditCmd.Flags().
		StringVar(&configEditSkipFilterItems, "skip-filter-untradable", "", "Skip filter untradable items (true/false)")
	configEditCmd.Flags().
		StringSliceVar(&configEditUpdateSecretKeys, "update-secret-keys", nil,
			"Keys of secrets to update")
	configEditCmd.Flags().
		StringSliceVar(&configEditUpdateSecretValues, "update-secret-values", nil,
			"Values of secrets to update (must match the keys in --update-secret-keys)")
	configEditCmd.MarkFlagsRequiredTogether("update-secret-keys", "update-secret-values")
}

const tabwriterPadding = 2

func printConfigAsTable(extended bool) error {
	if cfg == nil {
		return errors.New("configuration is not loaded")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, tabwriterPadding, ' ', 0)

	if !extended {
		return printBasicConfigTable(w)
	}

	return printExtendedConfigTable(w)
}

func printBasicConfigTable(w *tabwriter.Writer) error {
	_, err := fmt.Fprintln(
		w,
		"Project Name\tCreated At\tSteam ID 64\tSkip Filter Untradable Items",
	)
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	_, err = fmt.Fprintln(
		w,
		"------------\t----------\t------------\t--------------------------",
	)
	if err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	_, err = fmt.Fprintf(w, "%s\t%s\t%d\t%v\n",
		cfg.ProjectName,
		cfg.CreatedAt.Format(time.RFC3339),
		cfg.SteamID64,
		cfg.SkipFilterUntradableItems,
	)
	if err != nil {
		return fmt.Errorf("failed to write config values: %w", err)
	}

	return w.Flush()
}

func printExtendedConfigTable(w *tabwriter.Writer) error {
	_, err := fmt.Fprintln(
		w,
		"Project Name\tCreated At\tUpdated At\tCooldown Duration\tSkip Steam Services Check\tSteam ID 64\tSkip Steam User Check\tSkip Filter Untradable Items\tAdditional Items File",
	)
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	_, err = fmt.Fprintln(
		w,
		"------------\t----------\t----------\t------------------\t---------------------------\t------------\t----------------------\t--------------------------\t----------------------",
	)
	if err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\t%d\t%v\t%v\t%s\n",
		cfg.ProjectName,
		cfg.CreatedAt.Format(time.RFC3339),
		cfg.UpdatedAt.Format(time.RFC3339),
		cfg.CooldownDuration,
		cfg.SkipSteamServicesCheck,
		cfg.SteamID64,
		cfg.SkipSteamUserCheck,
		cfg.SkipFilterUntradableItems,
		cfg.AdditionalItemsFile,
	)
	if err != nil {
		return fmt.Errorf("failed to write config values: %w", err)
	}

	return w.Flush()
}

func getSecret(key secret.Key) (string, error) {
	value, err := secret.Load(key)
	if err != nil {
		return "", fmt.Errorf("failed to load secret %q: %w", key, err)
	}

	return value, nil
}

func parseBool(s string) (bool, error) {
	if s == "" {
		return false, errors.New("missing bool value")
	}
	return strconv.ParseBool(s)
}

func parseUint64(s string) (uint64, error) {
	if s == "" {
		return 0, errors.New("missing uint64 value")
	}
	return strconv.ParseUint(s, 10, 64)
}

const (
	timeUnitDay  = 24 * time.Hour
	timeUnitWeek = 7 * timeUnitDay
)

func parseExtendedDuration(s string) (time.Duration, error) {
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	unitMap := map[string]time.Duration{
		"d": timeUnitDay,
		"w": timeUnitWeek,
	}

	for unit, dur := range unitMap {
		if strings.HasSuffix(s, unit) {
			numStr := strings.TrimSuffix(s, unit)
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, err
			}
			return time.Duration(float64(dur) * num), nil
		}
	}

	return 0, fmt.Errorf("unknown duration format: %s", s)
}
