package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/devusSs/dropawp/internal/config"
	"github.com/devusSs/dropawp/internal/log"
	"github.com/devusSs/dropawp/internal/projects"
	"github.com/devusSs/dropawp/internal/stdin"
	"github.com/devusSs/dropawp/internal/storage"
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

		baseDir := filepath.Join(homeDir, ".dropawp")

		projs, err := projects.Load(baseDir)
		cobra.CheckErr(err)

		for _, project := range projs {
			if project.Name == name {
				cobra.CheckErr("Project with this name already exists")
			}
		}

		logsDir := filepath.Join(baseDir, "logs", name)
		err = log.Setup(logsDir)
		cobra.CheckErr(err)

		configFile := filepath.Join(baseDir, "configs", name+".json")
		storageDir := filepath.Join(baseDir, "storages", name)

		log.Debug("Creating new project",
			slog.String("name", name),
			slog.String("logsDir", logsDir),
			slog.String("configFile", configFile),
			slog.String("storageDir", storageDir),
		)

		conf, err := readConfigFromStdin(name)
		cobra.CheckErr(err)

		err = conf.Write(configFile)
		cobra.CheckErr(err)

		err = storage.CreateDir(storageDir)
		cobra.CheckErr(err)

		err = projects.SaveProject(baseDir, &projects.Project{
			Name:       name,
			LogsDir:    logsDir,
			ConfigFile: configFile,
			StorageDir: storageDir,
		})
		cobra.CheckErr(err)

		log.Info("Created new tracking project",
			slog.String("name", name),
			slog.String("logsDir", logsDir),
			slog.String("configFile", configFile),
			slog.String("storageDir", storageDir),
		)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func readConfigFromStdin(projectName string) (*config.Config, error) {
	logKeep, err := stdin.ReadFromStdin("Enter your log keep interval (default: 7d)")
	if err != nil {
		return nil, fmt.Errorf("failed to read log keep interval: %w", err)
	}

	if logKeep == "" {
		logKeep = "7d"
	}

	storageKeep, err := stdin.ReadFromStdin("Enter your storage keep interval (default: 30d)")
	if err != nil {
		return nil, fmt.Errorf("failed to read storage keep interval: %w", err)
	}

	if storageKeep == "" {
		storageKeep = "30d"
	}

	wantSteamCheck, err := stdin.ReadFromStdin("Do you want to enable Steam check? (default: true)")
	if err != nil {
		return nil, fmt.Errorf("failed to read Steam check preference: %w", err)
	}

	if wantSteamCheck == "" {
		wantSteamCheck = "true"
	}

	username, err := stdin.ReadFromStdin("Enter your username (default: project name)")
	if err != nil {
		return nil, fmt.Errorf("failed to read username: %w", err)
	}

	if username == "" {
		username = projectName
	}

	steamID64, err := stdin.ReadFromStdin("Enter your SteamID64 (required)")
	if err != nil {
		return nil, fmt.Errorf("failed to read SteamID64: %w", err)
	}

	if steamID64 == "" {
		return nil, errors.New("SteamID64 cannot be empty")
	}

	c, err := config.NewConfig(logKeep, storageKeep, wantSteamCheck, username, steamID64)
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	return c, nil
}
