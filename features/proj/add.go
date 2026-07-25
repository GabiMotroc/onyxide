package proj

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func projAdd(cmd *cobra.Command, args []string) error {
	return AddProject(args[0], args[1])
}

func AddProject(app string, location string) error {
	items, err := LoadProjects()
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("error loading projects: %w", err)
	}

	absoluteLocation, err := getLocation(app, location)
	if err != nil {
		return err
	}

	items = append(items, Project{AppType: app, Location: absoluteLocation})

	err = SaveProjects(items)

	if err != nil {
		return fmt.Errorf("error saving projects: %w", err)
	}
	return nil
}

func getLocation(app, location string) (string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	if len(location) == 0 {
		return currentPath, nil
	}

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
