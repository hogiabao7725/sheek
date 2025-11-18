package components

import (
	"strings"

	"sheek/internal/tui/styles"
)

// HighlightMatches highlights matching text in the command text based on search input.
// isSelected indicates if the item is currently selected, to use appropriate highlight style.
func HighlightMatches(text, searchInput string, isSelected bool) string {
	if strings.TrimSpace(searchInput) == "" {
		return text
	}

	searchLower := strings.ToLower(searchInput)
	textLower := strings.ToLower(text)

	// Find all match positions (case-insensitive)
	var matches []int
	start := 0
	for {
		idx := strings.Index(textLower[start:], searchLower)
		if idx == -1 {
			break
		}
		actualIdx := start + idx
		matches = append(matches, actualIdx)
		start = actualIdx + len(searchLower)
	}

	if len(matches) == 0 {
		return text
	}

	// Choose appropriate highlight style based on selection state
	style := styles.HighlightStyle
	if isSelected {
		style = styles.HighlightSelectedStyle
	}

	// Build highlighted text by inserting highlight styles at match positions
	var result strings.Builder
	lastPos := 0

	for _, matchStart := range matches {
		matchEnd := matchStart + len(searchLower)

		// Add text before match
		if matchStart > lastPos {
			result.WriteString(text[lastPos:matchStart])
		}

		// Extract the actual matched substring (preserving original case)
		matchedText := text[matchStart:matchEnd]

		// Highlight the matched text with appropriate style
		// For selected items, don't set background to preserve selected background
		if isSelected {
			// Use a style that only modifies foreground and underline, not background
			highlighted := styles.HighlightSelectedStyle.Copy().
				UnsetBackground().
				Render(matchedText)
			result.WriteString(highlighted)
		} else {
			highlighted := style.Render(matchedText)
			result.WriteString(highlighted)
		}

		lastPos = matchEnd
	}

	// Add remaining text after last match
	if lastPos < len(text) {
		result.WriteString(text[lastPos:])
	}

	return result.String()
}

