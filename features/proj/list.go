package proj

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func projList(cmd *cobra.Command, args []string) error {
	items, err := LoadProjects()
	if err != nil {
		return fmt.Errorf("error loading apps: %w", err)
	}

	w := new(tabwriter.Writer)
	w.Init(cmd.OutOrStdout(), 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tLOCATION")
	for _, item := range items {
		fmt.Fprintf(w, "%s\t%s\n", item.AppType, item.Location)
	}
	return w.Flush()
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  "Display all registered projects.",
	RunE:  projList,
}
