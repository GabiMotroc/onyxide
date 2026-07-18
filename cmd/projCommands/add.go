package projCommands

import (
	"OpenCli/data"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func projAdd(cmd *cobra.Command, args []string) error {
	return AddProject(args[0], args[1])
}

func AddProject(app string, location string) error {
	items, err := data.LoadProjects()
	if err != nil {
		return fmt.Errorf("error loading projects: %v", err)
	}

	dir, err := os.Getwd()
	fmt.Println(dir)
	items = append(items, data.Project{AppType: app, Location: location})

	err = data.SaveProjects(items)

	if err != nil {
		return fmt.Errorf("error saving apps: %v", err)
	}
	return nil
}

var addCmd = &cobra.Command{
	Use:   "add <command> <location>",
	Short: "Add a new project, if no location is provided, the current directory will be used",
	Long:  "Add a project with the given name.",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  projAdd,
}

func init() {
	ProjCmd.AddCommand(addCmd)
}
