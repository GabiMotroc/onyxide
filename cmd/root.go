package cmd

import (
	"onyxide/features/app"
	"onyxide/features/proj"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "o",
	Short: "onyxide - manage apps and projects",
	Long:  `onyxide is a CLI tool for managing apps and projects.`,
	Args:  cobra.ExactArgs(1),
	RunE:  proj.Open,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(app.Command)
	RootCmd.AddCommand(proj.Command)
	RootCmd.AddCommand(HookCmd)
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(UninitCmd)
}
