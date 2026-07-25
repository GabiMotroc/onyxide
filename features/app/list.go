package app

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func appList(cmd *cobra.Command, args []string) error {
	items, err := LoadApps()
	if err != nil {
		return fmt.Errorf("error loading apps: %w", err)
	}

	w := new(tabwriter.Writer)
	w.Init(cmd.OutOrStdout(), 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "NAME")
	for _, item := range items {
		fmt.Fprintf(w, "%s\n", item.Name)
	}
	return w.Flush()
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all apps",
	Long:  "Display all registered apps with their name and location.",
	RunE:  appList,
}
