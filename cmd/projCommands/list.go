package projCommands

import (
	"onyxide/data"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func projList(cmd *cobra.Command, args []string) {
	items, err := data.LoadProjects()
	if err != nil {
		_ = fmt.Errorf("error loading apps: %v", err)
	}

	w := new(tabwriter.Writer)
	w.Init(cmd.OutOrStdout(), 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tLOCATION")
	for _, item := range items {
		fmt.Fprintf(w, "%s\t%s\n", item.AppType, item.Location)
	}
	w.Flush()
}

var projListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  "Display all registered projects.",
	Run:   projList,
}

func init() {
	ProjCmd.AddCommand(projListCmd)
}
