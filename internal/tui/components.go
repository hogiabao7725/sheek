package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	PrimaryColor   = lipgloss.Color("#7D56F4")
	SecondaryColor = lipgloss.Color("#04B575")
	MutedColor     = lipgloss.Color("#626262")
	BorderColor    = lipgloss.Color("#FFFFFF")
	TextColor      = lipgloss.Color("#FFFFFF")
)

// Search Component Styles
var (
	// PromptStyle Prompt symbol (>)
	PromptStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	// Left panel - search input box
	SearchInputBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(BorderColor).
				Foreground(TextColor).
				PaddingLeft(1).
				PaddingRight(1).
				Height(1)

	// Right panel - mode badge
	ModeBadgeStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor).
			Foreground(SecondaryColor).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1).
			Height(1).
			Align(lipgloss.Center)

	// Container for search bar
	SearchContainerStyle = lipgloss.NewStyle().
				MarginLeft(1).
				MarginRight(1).
				MarginTop(1).
				MarginBottom(1)
)

// RenderSearchComponent renders the two-column search bar
func RenderSearchComponent(prompt string, inputValue string, mode string, terminalWidth int) string {
	// Calculate widths
	usableWidth := terminalWidth - 4 // Account for margins

	// 85% for input, 15% for mode
	// Subtract 4 for borders (2 chars per panel)
	contentWidth := usableWidth - 2
	inputPanelWidth := int(float64(contentWidth) * 0.85)
	modePanelWidth := contentWidth - inputPanelWidth

	// Build left panel content
	promptText := PromptStyle.Render(prompt)
	leftContent := promptText + inputValue

	leftPanel := SearchInputBoxStyle.
		Width(inputPanelWidth).
		Render(leftContent)

	// Build right panel content
	modeText := fmt.Sprintf("[%s]", mode)
	rightPanel := ModeBadgeStyle.
		Width(modePanelWidth).
		Render(modeText)

	// Join panels horizontally
	searchBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	return SearchContainerStyle.Render(searchBar)
}
