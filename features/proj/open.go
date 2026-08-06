package proj

import (
	"fmt"
	"onyxide/features/app"
	"os"
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

	command, terminal := openProject(foundProj)
	if terminal {
		command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
		return command.Run()
	}
	return command.Start()

}

func openProject(project Project) (*exec.Cmd, bool) {
	cmd := executeCommand(project.AppType, project.Location)
	return cmd, app.IsTerminal(project.AppType)
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
