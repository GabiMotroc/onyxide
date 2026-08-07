package cmd

import (
	"onyxide/cmd/shell"
	"onyxide/features/app"
	"onyxide/features/proj"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "onyxide",
	Short: "A smarter way to open your projects",
	Long: `onyxide tracks your projects by usage and opens the right one
with a single command. The more you open a project, the higher its
score — so the project you use most is always matched first.`,
	Version: "0.1.0",
	Annotations: map[string]string{
		"author": "Motroc Gabi	motrocgabi@gmail.com",
		"url":    "https://github.com/GabiMotroc/onyxide",
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.SetHelpTemplate(helpTemplate)
	RootCmd.SetUsageTemplate(usageTemplate)
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	RootCmd.AddCommand(app.Command)
	RootCmd.AddCommand(proj.Command)
	RootCmd.AddCommand(proj.OpenCmd)
	RootCmd.AddCommand(shell.InitCmd)
	RootCmd.AddCommand(shell.UninitCmd)
}

const helpTemplate = `{{.CommandPath}}{{if .Root.Version}} {{.Root.Version}}{{end}}{{if .Root.Annotations}}
{{index .Root.Annotations "author"}}
{{index .Root.Annotations "url"}}{{end}}

{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
`

const usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} <COMMAND>{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

EXAMPLES:
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Options:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Options:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
