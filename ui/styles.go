package ui

import "charm.land/lipgloss/v2"

var (
	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Margin(1, 0)

	DialogTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("62")).
				MarginBottom(1)
)

func BoolIcon(b bool) string {
	if b {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("●") // green
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("○") // gray
}
