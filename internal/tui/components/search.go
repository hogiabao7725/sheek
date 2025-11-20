package components

import (
	"fmt"

	"sheek/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

const (
	// SearchPanelWidthRatio is the default width ratio for the search panel (0.0-1.0)
	SearchPanelWidthRatio = 0.85
)

// RenderSearchComponent renders a two-column search bar with input and mode badge
func RenderSearchComponent(prompt, inputValue, mode string, terminalWidth int, horizontalMargin int) string {
	usableWidth := terminalWidth - (horizontalMargin * 2)
	contentWidth := usableWidth - 4

	inputWidth := int(float64(contentWidth) * SearchPanelWidthRatio)
	modeWidth := contentWidth - inputWidth

	leftContent := styles.PromptStyle.Render(prompt) + inputValue
	leftPanel := styles.SearchInputStyle.Width(inputWidth).Render(leftContent)

	modeText := fmt.Sprintf("[%s]", mode)
	rightPanel := styles.ModeBadgeStyle.Width(modeWidth).Render(modeText)

	searchBar := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	return styles.SearchContainerStyle.Render(searchBar)
}

