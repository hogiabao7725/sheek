package components

import (
	"fmt"
	"strings"

	"sheek/internal/history"
	"sheek/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// SearchMode represents the type of search to perform
type SearchMode string

const (
	// SearchModeExact performs exact substring matching
	SearchModeExact SearchMode = "Exact"
	// SearchModeFuzzy performs fuzzy matching
	SearchModeFuzzy SearchMode = "Fuzzy"
)

// RenderListComponent renders a sliding window of command list items with scrollbar
func RenderListComponent(commands []history.Command, fuzzyPositions map[int][]int, selectedIndex, terminalWidth, terminalHeight int, searchInput string, searchMode SearchMode) string {
	if len(commands) == 0 {
		emptyMessage := "No commands found"
		return styles.ListContainerStyle.
			Width(terminalWidth - (styles.HorizontalMargin * 2) - 2).
			Height(styles.ListContainerHeight).
			Render(styles.EmptyStateStyle.Render(emptyMessage))
	}

	startIndex, endIndex := calculateVisibleRange(len(commands), selectedIndex)

	// Calculate container width
	containerWidth := terminalWidth - (styles.HorizontalMargin * 2) - 2

	// Create scrollbar
	scrollbar := RenderScrollbar(len(commands), styles.MaxVisibleItems, selectedIndex, styles.ListContainerHeight)

	// Calculate item width based on whether scrollbar exists
	itemWidth := calculateItemWidth(containerWidth, scrollbar != "")

	// Render items with correct width for selected item highlighting
	items := renderCommandItems(commands, fuzzyPositions, startIndex, endIndex, selectedIndex, searchInput, searchMode, itemWidth)
	listContent := strings.Join(items, "\n")

	// If no scrollbar needed, return just the list
	if scrollbar == "" {
		return styles.ListContainerStyle.
			Width(containerWidth).
			Height(styles.ListContainerHeight).
			Render(listContent)
	}

	// Calculate available width for content (subtract scrollbar width and margin)
	contentWidth := containerWidth - 1 - 1 - 2 // -2 for padding

	// Create content with proper width (no border, just content)
	contentStyled := lipgloss.NewStyle().
		Width(contentWidth).
		Height(styles.ListContainerHeight - 2). // -2 for border
		Render(listContent)

	// Combine content and scrollbar horizontally
	combinedContent := lipgloss.JoinHorizontal(lipgloss.Top, contentStyled, scrollbar)

	// Wrap the combined content in a container with border
	return styles.ListContainerStyle.
		Width(containerWidth).
		Height(styles.ListContainerHeight).
		Render(combinedContent)
}

// calculateItemWidth determines the width of list items based on container and scrollbar presence
func calculateItemWidth(containerWidth int, hasScrollbar bool) int {
	if !hasScrollbar {
		// No scrollbar: container width - padding (1 on each side = 2 total)
		return containerWidth - 2
	}
	// With scrollbar: content width (container - border - padding - scrollbar space)
	return containerWidth - 1 - 1 - 2 // -2 for padding
}

// renderCommandItems creates styled items for the visible range with highlighting
func renderCommandItems(commands []history.Command, fuzzyPositions map[int][]int, start, end, selectedIndex int, searchInput string, searchMode SearchMode, itemWidth int) []string {
	items := make([]string, 0, styles.MaxVisibleItems)

	for i := start; i < end && i < len(commands); i++ {
		cmd := commands[i]
		isSelected := i == selectedIndex

		itemNumber := styles.ItemNumberStyle.Render(fmt.Sprintf("%d", cmd.Index))

		// Highlight matching text based on search mode
		var highlightedText string
		if searchMode == SearchModeFuzzy {
			// Use fuzzy highlighting with match positions
			matchPositions := fuzzyPositions[cmd.Index]
			highlightedText = HighlightFuzzyMatches(cmd.Text, matchPositions, isSelected)
		} else {
			// Use exact substring highlighting
			highlightedText = HighlightMatches(cmd.Text, searchInput, isSelected)
		}
		commandText := styles.CommandTextStyle.Render(highlightedText)

		itemContent := lipgloss.JoinHorizontal(lipgloss.Top, itemNumber, commandText)

		// Apply selected or normal style with full width to ensure background covers entire line
		var styledItem string
		if isSelected {
			// For selected items, ensure background covers full width
			// by rendering the style with explicit width
			styledItem = styles.ListItemSelectedStyle.
				Width(itemWidth).
				Render(itemContent)
		} else {
			styledItem = styles.ListItemStyle.
				Width(itemWidth).
				Render(itemContent)
		}
		items = append(items, styledItem)
	}

	return items
}

// calculateVisibleRange determines which items to display in a sliding window
// centered around the selected index. It returns the start (inclusive) and end (exclusive) indices.
func calculateVisibleRange(total, selectedIndex int) (start, end int) {
	if total <= styles.MaxVisibleItems {
		return 0, total
	}

	halfWindow := styles.MaxVisibleItems / 2

	// If selected item is near the beginning, show first MaxVisibleItems
	if selectedIndex < halfWindow {
		return 0, styles.MaxVisibleItems
	}

	// If selected item is near the end, show last MaxVisibleItems
	if selectedIndex > total-halfWindow {
		return total - styles.MaxVisibleItems, total
	}

	// Otherwise, center the window around the selected item
	return selectedIndex - halfWindow, selectedIndex + halfWindow
}

