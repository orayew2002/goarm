package utils

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF")).
			Bold(true).
			PaddingTop(1).
			PaddingBottom(1).
			Align(lipgloss.Center)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00BFFF")).
			Padding(0, 1).
			Width(40).
			Align(lipgloss.Left)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#999999")).
			MarginTop(1).
			Align(lipgloss.Center)

	sdfItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("#00BFFF"))

	sdfSelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color("10")).
				Bold(true)
)
