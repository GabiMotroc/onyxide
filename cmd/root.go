/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"OpenCli/cmd/appCommands"
	"OpenCli/cmd/projCommands"
)

func open(cmd *cobra.Command, args []string) {
	var command *exec.Cmd

	if len(args) == 0 {
		return
	}

	if isWindows() {
		command = exec.Command("cmd", "/c", "echo", args[0])
	} else {
		command = exec.Command("echo", args[0])
	}

	stdout, err := command.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "o",
	Short: "OpenCli - manage apps and projects",
	Long:  `OpenCli is a CLI tool for managing apps and projects.`,
	Args:  cobra.ArbitraryArgs,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: open,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func init() {
	RootCmd.AddCommand(appCommands.AppCmd)
	RootCmd.AddCommand(projCommands.ProjCmd)

	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
