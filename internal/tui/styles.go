package tui

import "github.com/charmbracelet/lipgloss"

// Color constants
const (
	primaryColorHex    = "#7D56F4"
	secondaryColorHex  = "#04B575"
	mutedColorHex      = "#626262"
	borderColorHex     = "#FFFFFF"
	textColorHex       = "#FFFFFF"
	backgroundColorHex = "#1A1A1A"
	selectedColorHex   = "#2D2D2D"
	accentColorHex     = "#FFD700"
)

var (
	primaryColor    = lipgloss.Color(primaryColorHex)
	secondaryColor  = lipgloss.Color(secondaryColorHex)
	mutedColor      = lipgloss.Color(mutedColorHex)
	borderColor     = lipgloss.Color(borderColorHex)
	textColor       = lipgloss.Color(textColorHex)
	backgroundColor = lipgloss.Color(backgroundColorHex)
	selectedColor   = lipgloss.Color(selectedColorHex)
	accentColor     = lipgloss.Color(accentColorHex)
)

// Style constants
const (
	searchPanelWidthRatio = 0.85
	maxVisibleItems       = 10
	listContainerHeight   = 12
	horizontalMargin      = 1
)

// Component styles
var (
	promptStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	searchInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(borderColor).
				Foreground(textColor).
				PaddingLeft(1).
				PaddingRight(1).
				Height(1)

	modeBadgeStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Foreground(secondaryColor).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1).
			Height(1).
			Align(lipgloss.Center)

	searchContainerStyle = lipgloss.NewStyle().
				MarginLeft(1).
				MarginRight(1).
				MarginTop(1).
				MarginBottom(1)

	listContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(borderColor).
				MarginLeft(1).
				MarginRight(1).
				MarginTop(1).
				Padding(1)

	listItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(1).
			PaddingRight(1)

	listItemSelectedStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(selectedColor).
				PaddingLeft(1).
				PaddingRight(1)

	itemNumberStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Width(4).
			Align(lipgloss.Right)

	commandTextStyle = lipgloss.NewStyle().
				Foreground(textColor).
				MarginLeft(1)

	emptyStateStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Align(lipgloss.Center).
			Padding(2)

	scrollbarTrackStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Width(1)

	scrollbarThumbStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Width(1)
)
