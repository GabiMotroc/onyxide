package appCommands

import (
	"onyxide/data"

	"github.com/spf13/cobra"
)

func clearApps(c *cobra.Command, args []string) {
	err := data.SaveApps([]data.App{})
	if err != nil {
		return
	}
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear apps",
	Long:  `Clear all apps from the system. This command will remove all app data and configurations.`,
	Run:   clearApps,
}

func init() {
	AppCmd.AddCommand(clearCmd)
}
