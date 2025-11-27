package config

// ColorConfig represents color configuration for UI elements
type ColorConfig struct {
	Primary    string `json:"primary"`   // Primary color (default: "#7D56F4")
	Secondary  string `json:"secondary"` // Secondary color (default: "#04B575")
	Text       string `json:"text"`      // Text color (default: "#FFFFFF")
	Border     string `json:"border"`    // Border color (default: "#FFFFFF")
	Muted      string `json:"muted"`     // Muted/secondary text color (default: "#626262")
	Selected   string `json:"selected"`  // Selected item background (default: "#3A3A5C")
	Highlight  string `json:"highlight"` // Highlight/accent color (default: "#FFD700")
	Background string `json:"bg"`        // Background color (default: "#1A1A1A")
}

// Config represents the application configuration
type Config struct {
	// Layout
	MaxItems      int  `json:"max_items"`      // Maximum items to display (default: 10)
	Height        int  `json:"height"`         // List container height (default: 12)
	Margin        int  `json:"margin"`         // Horizontal margin (default: 1)
	ShowTimestamp bool `json:"show_timestamp"` // Display command timestamp column (default: true)

	// Display
	Reverse    bool   `json:"reverse"`    // Reverse display order (default: false)
	Mode       string `json:"mode"`       // Search mode: "exact" or "fuzzy" (default: "exact")
	Contextual bool   `json:"contextual"` // Enable context-aware ranking (default: true)

	// Input
	Limit       int    `json:"limit"`       // Input character limit (default: 128)
	Placeholder string `json:"placeholder"` // Search placeholder text (default: "Search History...")

	// List
	Title string `json:"title"` // List title (default: "Recent Commands")

	// Colors
	Colors ColorConfig `json:"colors"` // Color configuration
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		MaxItems:      10,
		Height:        12,
		Margin:        1,
		ShowTimestamp: true,
		Reverse:       false,
		Mode:          "exact",
		Contextual:    true,
		Limit:         128,
		Placeholder:   "Search History...",
		Title:         "Recent Commands",
		Colors: ColorConfig{
			Primary:    "#7D56F4",
			Secondary:  "#04B575",
			Text:       "#FFFFFF",
			Border:     "#FFFFFF",
			Muted:      "#626262",
			Selected:   "#3A3A5C",
			Highlight:  "#FFD700",
			Background: "#1A1A1A",
		},
	}
}
