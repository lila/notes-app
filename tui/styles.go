package tui

import "github.com/charmbracelet/lipgloss"

var (
	// General
	docStyle = lipgloss.NewStyle().
			Margin(1, 2).
			Padding(0, 1). // Add horizontal padding
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700")) // Gold // Gold

	// Title
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#00BFFF")). // DeepSkyBlue
			Padding(0, 1)

	// Help
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	// Date View
	dateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")). // Gold
			Bold(true).
			Margin(1, 0, 0, 0)

	// Concept View
	conceptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ADFF2F")). // GreenYellow
			Bold(true)

	// List items
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00BFFF")). // DeepSkyBlue
				Bold(true)
	itemStyle = lipgloss.NewStyle()
)
