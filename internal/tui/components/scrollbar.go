package components

import (
	"strings"

	"sheek/internal/tui/styles"
)

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
			trackLines[i] = styles.ScrollbarThumbStyle.Render("█")
		} else {
			trackLines[i] = styles.ScrollbarTrackStyle.Render("│")
		}
	}

	return strings.Join(trackLines, "\n")
}

