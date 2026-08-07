package shell

import (
	"fmt"

	"github.com/spf13/cobra"
)

var registry = map[string]Shell{}

func Register(s Shell) {
	registry[s.Name()] = s
}

func initApp(cmd *cobra.Command, args []string) error {
	s := registry[args[0]]
	if s == nil {
		return fmt.Errorf("unsupported shell: %q (available: zsh, pwsh)", args[0])
	}
	_, err := fmt.Fprint(cmd.OutOrStdout(), s.Init())
	return err
}

var InitCmd = &cobra.Command{
	Use:     "init <shell>",
	Short:   "Print shell code to add integration",
	Example: `eval "$(onyxide init zsh)"`,
	Args:    cobra.ExactArgs(1),
	RunE:    initApp,
}

var UninitCmd = &cobra.Command{
	Use:     "uninit <shell>",
	Short:   "Print shell code to remove integration",
	Example: `eval "$(onyxide uninit zsh)"`,
	Args:    cobra.ExactArgs(1),
	RunE:    uninitApp(),
}

func uninitApp() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		s := registry[args[0]]
		if s == nil {
			return fmt.Errorf("unsupported shell: %q (available: zsh, pwsh)", args[0])
		}

		_, err := fmt.Fprint(cmd.OutOrStdout(), s.Uninit())
		return err
	}
}
