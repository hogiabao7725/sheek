package styles

import (
	"sheek/internal/config"

	"github.com/charmbracelet/lipgloss"
)

// Color variables - updated by InitializeStyles
var (
	primaryColor    lipgloss.Color
	secondaryColor  lipgloss.Color
	mutedColor      lipgloss.Color
	borderColor     lipgloss.Color
	textColor       lipgloss.Color
	selectedColor   lipgloss.Color
	accentColor     lipgloss.Color
	highlightColor  lipgloss.Color
	MutedColor      lipgloss.Color // Exported for external use
)

// Component styles - initialized by InitializeStyles
var (
	PromptStyle            lipgloss.Style
	SearchInputStyle       lipgloss.Style
	SearchPlaceholderStyle lipgloss.Style
	ModeBadgeStyle         lipgloss.Style
	SearchContainerStyle   lipgloss.Style
	ListContainerStyle     lipgloss.Style
	ListItemStyle          lipgloss.Style
	ListItemSelectedStyle  lipgloss.Style
	HighlightStyle         lipgloss.Style
	HighlightSelectedStyle lipgloss.Style
	ItemNumberStyle        lipgloss.Style
	CommandTextStyle       lipgloss.Style
	EmptyStateStyle        lipgloss.Style
	ScrollbarTrackStyle    lipgloss.Style
	ScrollbarThumbStyle    lipgloss.Style
)

// InitializeStyles initializes all styles based on the provided config
func InitializeStyles(cfg *config.Config) {
	// Update color variables from config
	primaryColor = lipgloss.Color(cfg.Colors.Primary)
	secondaryColor = lipgloss.Color(cfg.Colors.Secondary)
	mutedColor = lipgloss.Color(cfg.Colors.Muted)
	borderColor = lipgloss.Color(cfg.Colors.Border)
	textColor = lipgloss.Color(cfg.Colors.Text)
	selectedColor = lipgloss.Color(cfg.Colors.Selected)
	accentColor = lipgloss.Color(cfg.Colors.Highlight)
	highlightColor = lipgloss.Color(cfg.Colors.Highlight)
	MutedColor = mutedColor

	// Initialize all component styles
	PromptStyle = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	SearchInputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Foreground(textColor).
		PaddingLeft(1).PaddingRight(1).Height(1)
	SearchPlaceholderStyle = lipgloss.NewStyle().Foreground(mutedColor).Faint(true)
	ModeBadgeStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Foreground(secondaryColor).
		Bold(true).PaddingLeft(1).PaddingRight(1).Height(1).
		Align(lipgloss.Center)
	SearchContainerStyle = lipgloss.NewStyle().MarginLeft(1).MarginRight(1).MarginTop(1)
	ListContainerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		MarginLeft(1).MarginRight(1).
		PaddingLeft(1).PaddingRight(1)
	ListItemStyle = lipgloss.NewStyle().Foreground(textColor).PaddingLeft(1).PaddingRight(1)
	ListItemSelectedStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(selectedColor).
		Bold(true).PaddingLeft(1).PaddingRight(1)
	HighlightStyle = lipgloss.NewStyle().Foreground(highlightColor).Bold(true).Underline(true)
	HighlightSelectedStyle = lipgloss.NewStyle().Foreground(highlightColor).Bold(true).Underline(true)
	ItemNumberStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).Width(4).Align(lipgloss.Right)
	CommandTextStyle = lipgloss.NewStyle().Foreground(textColor).MarginLeft(1)
	EmptyStateStyle = lipgloss.NewStyle().Foreground(mutedColor).Align(lipgloss.Center).Padding(2)
	ScrollbarTrackStyle = lipgloss.NewStyle().Foreground(mutedColor).Width(1)
	ScrollbarThumbStyle = lipgloss.NewStyle().Foreground(primaryColor).Width(1)
}
