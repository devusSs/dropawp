package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/devusSs/dropawp/internal/config"
	"github.com/devusSs/dropawp/internal/log"
	"github.com/devusSs/dropawp/internal/projects"
	"github.com/devusSs/dropawp/internal/stdin"
	"github.com/devusSs/dropawp/internal/storage"
	"github.com/spf13/cobra"
)

var (
	baseDir        string
	loadedProjects []*projects.Project
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Add, list and delete tracking projects",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		baseDir = filepath.Join(homeDir, ".dropawp")

		loadedProjects, err = projects.Load(baseDir)
		cobra.CheckErr(err)
	},
}

var addProjectCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new tracking project",
	Run: func(_ *cobra.Command, _ []string) {
		name, err := stdin.ReadFromStdin("Enter the name of the new tracking project (required)")
		cobra.CheckErr(err)

		if name == "" {
			cobra.CheckErr("Project name cannot be empty")
		}

		for _, project := range loadedProjects {
			if project.Name == name {
				cobra.CheckErr(fmt.Errorf("Project with name '%s' already exists", name))
			}
		}

		logsDir := filepath.Join(baseDir, "logs", name)
		err = log.Setup(logsDir)
		cobra.CheckErr(err)

		configFile := filepath.Join(baseDir, "configs", name+".json")
		storageDir := filepath.Join(baseDir, "storages", name)

		conf, err := readConfigFromStdin(name)
		cobra.CheckErr(err)

		err = conf.Write(configFile)
		cobra.CheckErr(err)

		err = storage.CreateDir(storageDir)
		cobra.CheckErr(err)

		err = projects.SaveProject(baseDir, &projects.Project{
			Name:       name,
			Timestamp:  time.Now(),
			LogsDir:    logsDir,
			ConfigFile: configFile,
			StorageDir: storageDir,
		})
		cobra.CheckErr(err)

		fmt.Printf("Project '%s' added successfully.\n", name)
	},
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

const listProjectsTabPadding = 2

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracking projects",
	Run: func(_ *cobra.Command, _ []string) {
		if len(loadedProjects) == 0 {
			fmt.Println("No tracking projects found.")
			fmt.Println("Run 'dropawp new' to create a new project.")
			return
		}

		if listProjectsJSON {
			for _, project := range loadedProjects {
				fmt.Println(project.JSON())
				fmt.Println()
			}
			return
		}

		if listProjectsText {
			for _, project := range loadedProjects {
				fmt.Println(project.String())
				fmt.Println()
			}
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, listProjectsTabPadding, ' ', 0)
		fmt.Fprintln(w, "NAME\tTIMESTAMP\tLOGS DIR\tCONFIG FILE\tSTORAGE DIR")

		for _, project := range loadedProjects {
			fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%s\t%s\n",
				project.Name,
				project.Timestamp.Format(time.DateTime),
				project.LogsDir,
				project.ConfigFile,
				project.StorageDir,
			)
		}

		err := w.Flush()
		cobra.CheckErr(err)
	},
}

var deleteProjectCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a tracking project by its name",
	Run: func(_ *cobra.Command, args []string) {
		if deleteProjectName == "" {
			cobra.CheckErr(errors.New("Project name is required"))
		}

		newProjects := make([]*projects.Project, 0)
		appended := 0

		for _, project := range loadedProjects {
			if project.Name != deleteProjectName {
				newProjects = append(newProjects, project)
				appended++
			}
		}

		if appended == len(loadedProjects) {
			cobra.CheckErr(fmt.Errorf("Project '%s' not found", deleteProjectName))
		}

		err := projects.SaveProjects(baseDir, newProjects)
		cobra.CheckErr(err)

		// TODO: remove the project's directories and files

		fmt.Printf("Project '%s' deleted successfully.\n", deleteProjectName)
	},
}

var (
	listProjectsJSON bool
	listProjectsText bool
)

var deleteProjectName string

func init() {
	rootCmd.AddCommand(projectsCmd)

	projectsCmd.AddCommand(addProjectCmd)

	projectsCmd.AddCommand(listProjectsCmd)

	listProjectsCmd.Flags().
		BoolVarP(&listProjectsJSON, "json", "j", false, "Output projects information in JSON format")
	listProjectsCmd.Flags().
		BoolVarP(&listProjectsText, "text", "t", false, "Output projects information in text format")

	projectsCmd.AddCommand(deleteProjectCmd)
	deleteProjectCmd.Flags().
		StringVarP(&deleteProjectName, "name", "n", "", "Name of the project to delete")
	deleteProjectCmd.MarkFlagRequired("name")
}
