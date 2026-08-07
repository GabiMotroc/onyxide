package proj

import (
	"fmt"
	"onyxide/features/app"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var silent bool

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

func AddProject(appName string, workDir string, location string) error {
	apps, err := app.LoadApps()
	if err != nil {
		return fmt.Errorf("loading apps: %w", err)
	}
	if !app.ContainsAppName(apps, appName) {
		if silent {
			return nil
		}
		return fmt.Errorf("app %q is not registered", appName)
	}

	items, err := LoadProjects()

	if err != nil {
		return fmt.Errorf("error loading projects: %w", err)
	}

	absoluteLocation, err := getLocation(workDir, location)
	if err != nil {
		return err
	}

	items, err = appendProject(items, appName, absoluteLocation)
	if err != nil {
		return err
	}

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
	Use:   "add <command> [path]",
	Short: "Add a project",
	Long:  `Add a project for the given command. If no path is given, the current directory is used.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE:  projAdd,
}

func init() {
	addCmd.Flags().BoolVarP(&silent, "silent", "s", false, "suppress all output (used by the shell hook)")
}
