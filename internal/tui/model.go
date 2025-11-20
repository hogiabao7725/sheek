package tui

import (
	"sheek/internal/config"
	"sheek/internal/history"
	"sheek/internal/tui/components"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

// Model represents the application state
type Model struct {
	Input            textinput.Model
	List             list.Model
	Commands         []history.Command
	FilteredCommands []history.Command
	FuzzyPositions   map[int][]int // Map command index -> match positions for fuzzy highlighting
	SearchMode       SearchMode
	Width            int
	Height           int
	SelectedCommand  string // Command selected when user presses Enter
	LinesRendered    int    // Number of lines rendered by the UI (for cleanup)
	Placeholder      string
	Config           *config.Config // Application configuration
}

// NewModel creates a new Model with the given commands, config, and optional initial query
func NewModel(commands []history.Command, cfg *config.Config, initialQuery string) Model {
	// Reverse display order if configured
	if cfg.Reverse {
		reversed := make([]history.Command, len(commands))
		for i := 0; i < len(commands); i++ {
			reversed[i] = commands[len(commands)-1-i]
		}
		commands = reversed
	}

	items := components.CommandsToListItems(commands)

	in := textinput.New()
	in.Placeholder = ""
	in.Prompt = ""
	in.Focus()
	in.CharLimit = cfg.Limit
	in.SetValue(initialQuery)
	in.CursorEnd()

	l := list.New(items, list.NewDefaultDelegate(), 0, cfg.MaxItems)
	l.Title = cfg.Title
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	// Determine initial search mode from config
	initialSearchMode := SearchModeExact
	if cfg.Mode == "fuzzy" {
		initialSearchMode = SearchModeFuzzy
	}

	model := Model{
		Input:            in,
		List:             l,
		Commands:         commands,
		FilteredCommands: commands,
		SearchMode:       initialSearchMode,
		Placeholder:      cfg.Placeholder,
		Config:           cfg,
	}

	model = updateSearchResults(model)

	return model
}
