package projCommands

import (
	"OpenCli/data"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func projAdd(cmd *cobra.Command, args []string) error {
	items, err := data.LoadProjects()
	if err != nil {
		if !strings.Contains(err.Error(), "The system cannot find the file specified") {
			return fmt.Errorf("error loading projects: %v", err)
		}
	}

	location, err := getLocation(args)
	if err != nil {
		return err
	}

	fmt.Println(location)
	items = append(items, data.Project{AppType: args[0], Location: location})

	err = data.SaveProjects(items)

	if err != nil {
		return fmt.Errorf("error saving apps: %v", err)
	}
	return nil
}

func getLocation(args []string) (string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}

	if len(args) == 1 {
		return currentPath, nil
	}

	location := args[1]
	if filepath.IsAbs(location) {
		return location, nil
	}

	return filepath.Join(currentPath, location), nil
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
