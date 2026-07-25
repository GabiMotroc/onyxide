package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func runRemove(cmd *cobra.Command, args []string) error {
	apps, err := LoadApps()
	if err != nil {
		return fmt.Errorf("loading apps: %w", err)
	}

	filtered, found := RemoveByName(apps, args[0])
	if !found {
		return fmt.Errorf("app %q not found", args[0])
	}

	if err := SaveApps(filtered); err != nil {
		return fmt.Errorf("saving apps: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "app %q removed\n", args[0])
	return nil

}

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove an app",
	Long:  "Remove a registered app by name.",
	Args:  cobra.ExactArgs(1),
	RunE:  runRemove,
}
