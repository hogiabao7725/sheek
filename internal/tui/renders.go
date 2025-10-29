package tui

import (
	"fmt"
	"strings"

	"sheek/internal/history"

	"github.com/charmbracelet/lipgloss"
)

// RenderSearchComponent renders a two-column search bar with input and mode badge
func RenderSearchComponent(prompt, inputValue, mode string, terminalWidth int) string {

	usableWidth := terminalWidth - (horizontalMargin * 2)
	contentWidth := usableWidth - 4

	inputWidth := int(float64(contentWidth) * searchPanelWidthRatio)
	modeWidth := contentWidth - inputWidth

	leftContent := promptStyle.Render(prompt) + inputValue
	leftPanel := searchInputStyle.Width(inputWidth).Render(leftContent)

	modeText := fmt.Sprintf("[%s]", mode)
	rightPanel := modeBadgeStyle.Width(modeWidth).Render(modeText)

	searchBar := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	return searchContainerStyle.Render(searchBar)
}

// RenderListComponent renders a sliding window of command list items with scrollbar
func RenderListComponent(commands []history.Command, selectedIndex, terminalWidth, terminalHeight int) string {
	if len(commands) == 0 {
		emptyMessage := "No commands found"
		return listContainerStyle.
			Width(terminalWidth - (horizontalMargin * 2) - 2).
			Height(listContainerHeight).
			Render(emptyStateStyle.Render(emptyMessage))
	}

	startIndex, endIndex := calculateVisibleRange(len(commands), selectedIndex)
	items := renderCommandItems(commands, startIndex, endIndex, selectedIndex)
	listContent := strings.Join(items, "\n")

	// Create scrollbar
	scrollbar := RenderScrollbar(len(commands), maxVisibleItems, selectedIndex, listContainerHeight)

	// Calculate container width
	containerWidth := terminalWidth - (horizontalMargin * 2) - 2

	// If no scrollbar needed, return just the list
	if scrollbar == "" {
		return listContainerStyle.
			Width(containerWidth).
			Height(listContainerHeight).
			Render(listContent)
	}

	// Calculate available width for content (subtract scrollbar width and margin)
	contentWidth := containerWidth - 1 - 1 - 2 // -2 for padding

	// Create content with proper width (no border, just content)
	contentStyled := lipgloss.NewStyle().
		Width(contentWidth).
		Height(listContainerHeight - 2). // -2 for border
		Render(listContent)

	// Combine content and scrollbar horizontally
	combinedContent := lipgloss.JoinHorizontal(lipgloss.Top, contentStyled, scrollbar)

	// Wrap the combined content in a container with border
	return listContainerStyle.
		Width(containerWidth).
		Height(listContainerHeight).
		Render(combinedContent)
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

// RenderScrollbar creates a vertical scrollbar similar to fzf
func RenderScrollbar(totalItems, visibleItems, selectedIndex, height int) string {
	if totalItems <= visibleItems {
		return ""
	}

	// Calculate scrollbar position and size
	scrollbarHeight := height - 2 // Account for borders
	thumbHeight := max(1, int(float64(visibleItems)/float64(totalItems)*float64(scrollbarHeight)))

	// Calculate thumb position based on selected index
	scrollRatio := float64(selectedIndex) / float64(totalItems-1)
	thumbPosition := int(scrollRatio * float64(scrollbarHeight-thumbHeight))

	// Create scrollbar track
	trackLines := make([]string, scrollbarHeight)
	for i := 0; i < scrollbarHeight; i++ {
		if i >= thumbPosition && i < thumbPosition+thumbHeight {
			trackLines[i] = scrollbarThumbStyle.Render("█")
		} else {
			trackLines[i] = scrollbarTrackStyle.Render("│")
		}
	}

	return strings.Join(trackLines, "\n")
}
