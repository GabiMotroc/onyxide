package proj

import (
	"github.com/spf13/cobra"
)

func clearProjects(c *cobra.Command, args []string) error {
	return SaveProjects([]Project{})
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear projects",
	Long:  `Clear all projects from the system. This command will remove all projects data and configurations.`,
	RunE:  clearProjects,
}
