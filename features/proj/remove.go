package proj

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <location>",
	Short: "Remove project",
	Long:  "Remove a project",
	Args:  cobra.ExactArgs(1),
	RunE:  runRemove,
}

func runRemove(cmd *cobra.Command, args []string) error {
	projects, err := LoadProjects()
	if err != nil {
		return fmt.Errorf("loading projects: %w", err)
	}

	matches := AllMatchingProjects(args[0], projects)
	if len(matches) == 0 {
		return fmt.Errorf("project %q not found", args[0])
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%d project(s) match %q:\n", len(matches), args[0])
	for i, p := range matches {
		fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s (via %s)\n", i+1, p.Location, p.AppType)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Delete all %d project(s)? [y/N]: ", len(matches))
	answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading confirmation: %w", err)
	}
	if !strings.EqualFold(strings.TrimSpace(answer), "y") {
		fmt.Fprintln(cmd.OutOrStdout(), "aborted")
		return nil
	}

	if err := SaveProjects(RemoveMatching(args[0], projects)); err != nil {
		return fmt.Errorf("saving projects: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%d projects %q removed\n", len(matches), args[0])
	return nil

}
