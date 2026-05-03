package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	Cyan  = lipgloss.Color("86")
	Green = lipgloss.Color("42")
	Red   = lipgloss.Color("196")
	Gray  = lipgloss.Color("240")
	White = lipgloss.Color("252")

	ActiveTabBg = lipgloss.Color("62")
	SelectedBg  = lipgloss.Color("237")

	// Panels
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Gray).
			Padding(1, 2)

	// Typography
	TitleStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	// Tabs
	TabStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(Gray)

	ActiveTabStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(ActiveTabBg).
			Foreground(White).
			Bold(true)

	// Lists
	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
			Background(SelectedBg).
			Foreground(White).
			Bold(true).
			PaddingLeft(2)

	// Statuses
	StatusRunningStyle = lipgloss.NewStyle().Foreground(Green)
	StatusExitedStyle  = lipgloss.NewStyle().Foreground(Red)

	// Footer
	FooterStyle = lipgloss.NewStyle().
			Foreground(Gray).
			MarginTop(1)
)