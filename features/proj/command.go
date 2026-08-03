package proj

import "github.com/spf13/cobra"

var Command = &cobra.Command{
	Use:   "proj",
	Short: "Manage projects",
	Long:  "Create, list, and manage projects.",
	RunE:  startInteractive,
}

func init() {
	Command.AddCommand(addCmd, listCmd, clearCmd, removeCmd)
}
