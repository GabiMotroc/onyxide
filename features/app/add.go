package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func appAdd(cmd *cobra.Command, args []string) error {
	items, err := LoadApps()
	if err != nil {
		return fmt.Errorf("error loading apps: %w", err)
	}

	if ContainsAppName(items, args[0]) {
		return fmt.Errorf("app with name %s already exists", args[0])
	}
	items = append(items, App{Name: args[0]})

	err = SaveApps(items)

	if err != nil {
		return fmt.Errorf("error saving apps: %w", err)
	}
	return nil
}

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new app",
	Long:  "Add an app by name.",
	Args:  cobra.ExactArgs(1),
	RunE:  appAdd,
}
