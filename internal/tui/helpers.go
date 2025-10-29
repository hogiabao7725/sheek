package tui

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
