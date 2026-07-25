package proj

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

func Open(cmd *cobra.Command, args []string) error {
	projects, err := LoadProjects()
	if err != nil {
		return err
	}

	needle := strings.ToLower(args[0])

	found := false
	var foundProj Project
	for _, project := range projects {
		if strings.Contains(
			strings.ToLower(project.Location),
			needle,
		) {
			foundProj = project
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("no project matching %s", args[0])
	}

	fmt.Printf("opening %s using %s", foundProj.Location, foundProj.AppType)
	command := executeCommand(foundProj.AppType, foundProj.Location)

	return command.Start()
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
