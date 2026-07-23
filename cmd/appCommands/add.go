package appCommands

import (
	"fmt"
	"onyxide/data"

	"github.com/spf13/cobra"
)

func appAdd(cmd *cobra.Command, args []string) error {
	items, err := data.LoadApps()
	if err != nil {
		return fmt.Errorf("error loading apps: %v", err)
	}

	if data.ContainsAppName(items, args[0]) {
		return fmt.Errorf("app with name %s already exists", args[0])
	}
	items = append(items, data.App{Name: args[0]})

	err = data.SaveApps(items)

	if err != nil {
		return fmt.Errorf("error saving apps: %v", err)
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

func init() {
	AppCmd.AddCommand(addCmd)
}
