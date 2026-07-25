package cmd

import (
	"fmt"
	"onyxide/features/app"
	"onyxide/features/proj"
	"strings"

	"github.com/spf13/cobra"
)

var pwd, raw string

func hook(c *cobra.Command, args []string) error {
	fmt.Printf("received pwd=%s raw=%s\n", pwd, raw)

	apps, err := app.LoadApps()
	if err != nil {
		return err
	}

	split := strings.Split(raw, " ")

	if len(split) < 2 {
		fmt.Printf("received incompatible command")
		//return fmt.Errorf("invalid command")
		return nil
	}

	triggeredCmd := split[0]

	for _, a := range apps {
		if triggeredCmd == a.Name {
			err := proj.AddProject(triggeredCmd, split[1])
			if err != nil {
				return err
			}
			fmt.Printf("successfully saved pwd=%s raw=%s\n", pwd, raw)
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
	HookCmd.Flags().StringVar(&pwd, "pwd", "", "working directory")
	HookCmd.Flags().StringVar(&raw, "raw", "", "raw command")
}
