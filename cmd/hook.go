package cmd

import (
	"fmt"
	"onyxide/features/app"
	"onyxide/features/proj"
	"strings"

	"github.com/spf13/cobra"
)

var hookedLocation, hookedCmd string

func hook(c *cobra.Command, args []string) error {
	fmt.Printf("received pwd=%s raw=%s\n", hookedLocation, hookedCmd)

	apps, err := app.LoadApps()
	if err != nil {
		return err
	}

	split := strings.Split(hookedCmd, " ")

	if len(split) < 2 {
		fmt.Printf("received incompatible command")
		return nil
	}

	triggeredCmd := split[0]

	for _, a := range apps {
		if triggeredCmd == a.Name {
			err := proj.AddProject(triggeredCmd, hookedLocation, split[1])
			if err != nil {
				return err
			}
			fmt.Printf("successfully saved pwd=%s raw=%s\n", hookedLocation, hookedCmd)
		}
	}
	return nil
}

var HookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Hook",
	Args:  cobra.NoArgs,
	Long:  ``,
	RunE:  hook,
}

func init() {
	HookCmd.Flags().StringVar(&hookedLocation, "pwd", "", "working directory")
	HookCmd.Flags().StringVar(&hookedCmd, "raw", "", "raw command")
}
