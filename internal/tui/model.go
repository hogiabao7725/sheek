package tui

import (
	"sheek/internal/history"
	"sheek/internal/tui/components"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

// Model represents the application state
type Model struct {
	Input           textinput.Model
	List            list.Model
	Commands        []history.Command
	FilteredCommands []history.Command
	SearchMode      SearchMode
	Width           int
	Height          int
}

// NewModel creates a new Model with the given commands
func NewModel(commands []history.Command) Model {
	items := components.CommandsToListItems(commands)

	in := textinput.New()
	in.Placeholder = "Search History..."
	in.Focus()
	in.CharLimit = 128

	l := list.New(items, list.NewDefaultDelegate(), 0, 10)
	l.Title = "Recent Commands"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return Model{
		Input:           in,
		List:            l,
		Commands:        commands,
		FilteredCommands: commands,
		SearchMode:      SearchModeExact,
	}
}
