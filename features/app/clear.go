package app

import (
	"github.com/spf13/cobra"
)

func clearApps(c *cobra.Command, args []string) error {
	return SaveApps([]App{})
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear apps",
	Long:  `Clear all apps from the system. This command will remove all app data and configurations.`,
	RunE:  clearApps,
}
