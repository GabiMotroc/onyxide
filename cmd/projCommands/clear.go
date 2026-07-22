package projCommands

import (
	"onyxide/data"

	"github.com/spf13/cobra"
)

func clearProjects(c *cobra.Command, args []string) {
	err := data.SaveProjects([]data.Project{})
	if err != nil {
		return
	}
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear projects",
	Long:  `Clear all projects from the system. This command will remove all projects data and configurations.`,
	Run:   clearProjects,
}

func init() {
	ProjCmd.AddCommand(clearCmd)
}
