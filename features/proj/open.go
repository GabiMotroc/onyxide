package proj

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

func Open(cmd *cobra.Command, args []string) error {
	projects, err := LoadProjects()
	if err != nil {
		return err
	}

	found, foundProj := FirstMatchingProject(args[0], projects)

	if !found {
		return fmt.Errorf("no project matching %s", args[0])
	}

	fmt.Fprintf(cmd.OutOrStdout(), "opening %s using %s", foundProj.Location, foundProj.AppType)
	command := openProject(foundProj)

	return command.Start()
}

func openProject(project Project) *exec.Cmd {
	return executeCommand(project.AppType, project.Location)
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

func isWindows() bool {
	return runtime.GOOS == "windows"
}
