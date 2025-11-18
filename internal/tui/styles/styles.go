package styles

import "github.com/charmbracelet/lipgloss"

// Color hex constants
const (
	primaryColorHex    = "#7D56F4"
	secondaryColorHex  = "#04B575"
	mutedColorHex      = "#626262"
	borderColorHex     = "#FFFFFF"
	textColorHex       = "#FFFFFF"
	backgroundColorHex = "#1A1A1A"
	selectedColorHex   = "#3A3A5C"
	highlightColorHex  = "#FFD700"
	accentColorHex     = "#FFD700"
)

// Color variables for use in styles
var (
	// Exported colors
	MutedColor = lipgloss.Color(mutedColorHex)

	// Internal colors (used in styles)
	primaryColor    = lipgloss.Color(primaryColorHex)
	secondaryColor  = lipgloss.Color(secondaryColorHex)
	mutedColor      = MutedColor
	borderColor     = lipgloss.Color(borderColorHex)
	textColor       = lipgloss.Color(textColorHex)
	backgroundColor = lipgloss.Color(backgroundColorHex)
	selectedColor   = lipgloss.Color(selectedColorHex)
	accentColor     = lipgloss.Color(accentColorHex)
	highlightColor  = lipgloss.Color(highlightColorHex)
)

// Layout constants
const (
	SearchPanelWidthRatio = 0.85
	MaxVisibleItems       = 10
	ListContainerHeight   = 12
	HorizontalMargin      = 1
)

// Component styles - exported for use in components package
var (
	PromptStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	SearchInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(borderColor).
				Foreground(textColor).
				PaddingLeft(1).
				PaddingRight(1).
				Height(1)

	ModeBadgeStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Foreground(secondaryColor).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1).
			Height(1).
			Align(lipgloss.Center)

	SearchContainerStyle = lipgloss.NewStyle().
				MarginLeft(1).
				MarginRight(1).
				MarginTop(1).
				MarginBottom(1)

	ListContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(borderColor).
				MarginLeft(1).
				MarginRight(1).
				MarginTop(1).
				Padding(1)

	ListItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(1).
			PaddingRight(1)

	ListItemSelectedStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(selectedColor).
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(highlightColor).
			Bold(true).
			Underline(true)

	HighlightSelectedStyle = lipgloss.NewStyle().
			Foreground(highlightColor).
			Bold(true).
			Underline(true)

	ItemNumberStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Width(4).
			Align(lipgloss.Right)

	CommandTextStyle = lipgloss.NewStyle().
				Foreground(textColor).
				MarginLeft(1)

	EmptyStateStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Align(lipgloss.Center).
			Padding(2)

	ScrollbarTrackStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Width(1)

	ScrollbarThumbStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Width(1)
)

