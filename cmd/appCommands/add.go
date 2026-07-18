package appCommands

import (
	"OpenCli/data"
	"fmt"

	"github.com/spf13/cobra"
)

func appAdd(cmd *cobra.Command, args []string) error {
	items, err := data.LoadApps()
	if err != nil {
		return fmt.Errorf("error loading apps: %v", err)
	}

	items = append(items, data.App{Name: args[0]})

	err = data.SaveApps(items)

	if err != nil {
		return fmt.Errorf("error saving apps: %v", err)
	}
	return nil
}

var addCmd = &cobra.Command{
	Use:   "add <name> <location>",
	Short: "Add a new app",
	Long:  "Add an app by name and location.",
	Args:  cobra.ExactArgs(2),
	RunE:  appAdd,
}

func init() {
	AppCmd.AddCommand(addCmd)
}
