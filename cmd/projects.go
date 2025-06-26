package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/devusSs/dropawp/internal/projects"
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage and list tracking projects",
}

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracking projects",
	Run: func(_ *cobra.Command, _ []string) {
		baseDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		baseDir = filepath.Join(baseDir, ".dropawp")

		projects, err := projects.Load(baseDir)
		cobra.CheckErr(err)

		if len(projects) == 0 {
			fmt.Println("No tracking projects found.")
			fmt.Println("Run 'dropawp new' to create a new project.")
			return
		}

		if listProjectsJSON {
			for _, project := range projects {
				fmt.Println(project.JSON())
				fmt.Println()
			}
			return
		}

		if listProjectsText {
			for _, project := range projects {
				fmt.Println(project.String())
				fmt.Println()
			}
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tLOGS DIR\tCONFIG FILE\tSTORAGE DIR")

		for _, project := range projects {
			fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%s\n",
				project.Name,
				project.LogsDir,
				project.ConfigFile,
				project.StorageDir,
			)
		}

		err = w.Flush()
		cobra.CheckErr(err)
	},
}

var (
	listProjectsJSON bool
	listProjectsText bool
)

func init() {
	rootCmd.AddCommand(projectsCmd)

	projectsCmd.AddCommand(listProjectsCmd)

	listProjectsCmd.Flags().
		BoolVarP(&listProjectsJSON, "json", "j", false, "Output projects information in JSON format")
	listProjectsCmd.Flags().
		BoolVarP(&listProjectsText, "text", "t", false, "Output projects information in text format")
}
