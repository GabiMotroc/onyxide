package cmd

import (
	"onyxide/data"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"onyxide/cmd/appCommands"
	"onyxide/cmd/projCommands"
)

func open(cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		return fmt.Errorf("no arguments provided")
	}

	projects, err := data.LoadProjects()
	if err != nil {
		return err
	}

	var foundProj data.Project
	for _, project := range projects {
		if strings.Contains(
			strings.ToLower(project.Location),
			strings.ToLower(args[0]),
		) {
			fmt.Printf("opening %s using %s", project.Location, project.AppType)
			foundProj = project
			break
		}
	}

	command := executeCommand(foundProj.AppType, foundProj.Location)

	err = command.Start()

	if err != nil {
		return err
	}

	//fmt.Println(string(stdout))
	return nil
}

func executeCommand(app string, location string) *exec.Cmd {
	var command *exec.Cmd
	if isWindows() {
		command = exec.Command("cmd", "/c", app, location)
	} else {
		command = exec.Command(app, location)
	}
	return command
}

var RootCmd = &cobra.Command{
	Use:   "o",
	Short: "onyxide - manage apps and projects",
	Long:  `onyxide is a CLI tool for managing apps and projects.`,
	Args:  cobra.ExactArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: open,
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
	//RootCmd.AddCommand(HookCmd)
	//RootCmd.AddCommand(InitCmd)
	//RootCmd.AddCommand(UninitCmd)

	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
