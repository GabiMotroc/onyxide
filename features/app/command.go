package app

import (
	"github.com/spf13/cobra"
)

// Command is the parent "app" subcommand, exported so root.go can attach it.
var Command = &cobra.Command{
	Use:   "app",
	Short: "Manage apps",
	Long:  "Create, list, and manage applications.",
	RunE:  startInteractive,
}

func init() {
	Command.AddCommand(addCmd, listCmd, clearCmd, removeCmd)
}
