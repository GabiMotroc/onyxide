package proj

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func projAdd(cmd *cobra.Command, args []string) error {
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	location := ""
	if len(args) > 1 {
		location = args[1]
	}
	return AddProject(args[0], workDir, location)

}

func AddProject(app string, workDir string, location string) error {
	items, err := LoadProjects()

	if err != nil {
		return fmt.Errorf("error loading projects: %w", err)
	}

	absoluteLocation, err := getLocation(workDir, location)
	if err != nil {
		return err
	}

	if ContainsLocation(items, absoluteLocation) {
		return fmt.Errorf("project at %q already exists", absoluteLocation)
	}

	items = append(items, Project{AppType: app, Location: absoluteLocation})

	err = SaveProjects(items)

	if err != nil {
		return fmt.Errorf("error saving projects: %w", err)
	}
	return nil
}

func getLocation(workDir string, location string) (string, error) {
	if len(location) == 0 {
		return workDir, nil
	}

	if filepath.IsAbs(location) {
		return location, nil
	}

	return filepath.Join(workDir, location), nil
}

var addCmd = &cobra.Command{
	Use:   "add <command> <location>",
	Short: "Add a new project, if no location is provided, the current directory will be used",
	Long:  "Add a project with the given name.",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  projAdd,
}
