package tui

import (
	"fmt"
	"strings"

	"sheek/internal/history"

	"github.com/charmbracelet/lipgloss"
)

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
	horizontalMargin      = 2
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
)

// RenderSearchComponent renders a two-column search bar with input and mode badge
func RenderSearchComponent(prompt, inputValue, mode string, terminalWidth int) string {
	usableWidth := terminalWidth - (horizontalMargin * 2)
	contentWidth := usableWidth - 2

	inputWidth := int(float64(contentWidth) * searchPanelWidthRatio)
	modeWidth := contentWidth - inputWidth

	leftContent := promptStyle.Render(prompt) + inputValue
	leftPanel := searchInputStyle.Width(inputWidth).Render(leftContent)

	modeText := fmt.Sprintf("[%s]", mode)
	rightPanel := modeBadgeStyle.Width(modeWidth).Render(modeText)

	searchBar := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	return searchContainerStyle.Render(searchBar)
}

// RenderListComponent renders a sliding window of command list items
func RenderListComponent(commands []history.Command, selectedIndex, terminalWidth, terminalHeight int) string {
	if len(commands) == 0 {
		emptyMessage := "No commands found"
		return listContainerStyle.
			Width(terminalWidth - (horizontalMargin * 2)).
			Height(listContainerHeight).
			Render(emptyStateStyle.Render(emptyMessage))
	}

	startIndex, endIndex := calculateVisibleRange(len(commands), selectedIndex)
	items := renderCommandItems(commands, startIndex, endIndex, selectedIndex)
	listContent := strings.Join(items, "\n")

	return listContainerStyle.
		Width(terminalWidth - (horizontalMargin * 2)).
		Height(listContainerHeight).
		Render(listContent)
}

// calculateVisibleRange determines which items to display in a sliding window
func calculateVisibleRange(total, selectedIndex int) (start, end int) {
	if total <= maxVisibleItems {
		return 0, total
	}

	halfWindow := maxVisibleItems / 2

	if selectedIndex < halfWindow {
		return 0, maxVisibleItems
	}

	if selectedIndex > total-halfWindow {
		return total - maxVisibleItems, total
	}

	return selectedIndex - halfWindow, selectedIndex + halfWindow
}

// renderCommandItems creates styled items for the visible range
func renderCommandItems(commands []history.Command, start, end, selectedIndex int) []string {
	items := make([]string, 0, maxVisibleItems)

	for i := start; i < end && i < len(commands); i++ {
		cmd := commands[i]
		isSelected := i == selectedIndex

		itemNumber := itemNumberStyle.Render(fmt.Sprintf("%d", cmd.Index))
		commandText := commandTextStyle.Render(cmd.Text)
		itemContent := lipgloss.JoinHorizontal(lipgloss.Top, itemNumber, commandText)

		if isSelected {
			items = append(items, listItemSelectedStyle.Render(itemContent))
		} else {
			items = append(items, listItemStyle.Render(itemContent))
		}
	}

	return items
}
