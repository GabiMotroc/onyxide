package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func initApp(c *cobra.Command, args []string) error {
	switch args[0] {
	case "zsh":
		fmt.Print(zshInitScript)
		return nil
	}

	return errors.New("no shell found")
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize open",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE:  initApp,
}

func init() {
	RootCmd.AddCommand(InitCmd)
}

const zshInitScript = `
alias o='onyxide'

autoload -Uz add-zsh-hook

_mycli_preexec() {
  typeset -g MYCLI_LAST_CMD="$1"
  typeset -g MYCLI_LAST_PWD="$PWD"
}

_mycli_precmd() {
  if [[ -n "$MYCLI_LAST_CMD" ]]; then
    command onyxide hook --pwd "$MYCLI_LAST_PWD" --raw "$MYCLI_LAST_CMD"
    unset MYCLI_LAST_CMD
    unset MYCLI_LAST_PWD
  fi
}

add-zsh-hook preexec _mycli_preexec
add-zsh-hook precmd _mycli_precmd
`
