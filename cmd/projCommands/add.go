package projCommands

import (
	"OpenCli/data"
	"fmt"

	"github.com/spf13/cobra"
)

func projAdd(cmd *cobra.Command, args []string) error {
	items, err := data.LoadProjects()
	if err != nil {
		return fmt.Errorf("error loading projects: %v", err)
	}

	items = append(items, data.Project{AppType: args[0], Location: args[1]})

	err = data.SaveProjects(items)

	if err != nil {
		return fmt.Errorf("error saving apps: %v", err)
	}
	return nil
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new project",
	Long:  "Add a project with the given name.",
	Args:  cobra.MaximumNArgs(2),
	RunE:  projAdd,
}

func init() {
	ProjCmd.AddCommand(addCmd)
}
