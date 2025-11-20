package tui

// SearchMode represents the type of search to perform
type SearchMode string

const (
	// SearchModeExact performs exact substring matching
	SearchModeExact SearchMode = "Exact"
	// SearchModeFuzzy performs fuzzy matching
	SearchModeFuzzy SearchMode = "Fuzzy"
)

// String returns the string representation of SearchMode
func (s SearchMode) String() string {
	return string(s)
}

// Toggle switches between Exact and Fuzzy modes
func (s SearchMode) Toggle() SearchMode {
	if s == SearchModeExact {
		return SearchModeFuzzy
	}
	return SearchModeExact
}
