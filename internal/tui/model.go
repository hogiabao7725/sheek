package tui

import (
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
}

// NewModel creates a new Model with the given commands and optional initial query
func NewModel(commands []history.Command, initialQuery string) Model {
	items := components.CommandsToListItems(commands)

	in := textinput.New()
	in.Placeholder = ""
	in.Prompt = ""
	in.Focus()
	in.CharLimit = 128
	in.SetValue(initialQuery)
	in.CursorEnd()

	l := list.New(items, list.NewDefaultDelegate(), 0, 10)
	l.Title = "Recent Commands"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	model := Model{
		Input:            in,
		List:             l,
		Commands:         commands,
		FilteredCommands: commands,
		SearchMode:       SearchModeExact,
		Placeholder:      "Search History...",
	}

	model = updateSearchResults(model)

	return model
}
