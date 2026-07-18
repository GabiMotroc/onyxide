package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var UninitCmd = &cobra.Command{
	Use:   "uninit [shell]",
	Short: "Print shell code to remove integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "zsh":
			fmt.Print(zshUninitScript)
		}
		return nil
	},
}

const zshUninitScript = `
add-zsh-hook -d preexec _mycli_preexec
add-zsh-hook -d precmd _mycli_precmd
unfunction _mycli_preexec _mycli_precmd 2>/dev/null
unset MYCLI_LAST_CMD
`

func init() {
	RootCmd.AddCommand(UninitCmd)
}
