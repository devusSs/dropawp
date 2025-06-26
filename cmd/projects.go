package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/devusSs/dropawp/internal/projects"
	"github.com/spf13/cobra"
)

var (
	baseDir        string
	loadedProjects []*projects.Project
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage and list tracking projects",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		baseDir = filepath.Join(homeDir, ".dropawp")

		loadedProjects, err = projects.Load(baseDir)
		cobra.CheckErr(err)
	},
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
		fmt.Fprintln(w, "NAME\tLOGS DIR\tCONFIG FILE\tSTORAGE DIR")

		for _, project := range loadedProjects {
			fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%s\n",
				project.Name,
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
