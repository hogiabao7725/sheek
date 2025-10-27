package tui

import (
	"sheek/internal/history"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	Input      textinput.Model
	List       list.Model
	Commands   []history.Command
	SearchMode string
	MatchCount int
	Width      int
	Height     int
}

func NewModel(commands []history.Command) Model {
	items := CommandsToListItems(commands)

	in := textinput.New()
	in.Placeholder = "Search History..."
	in.Focus()
	in.CharLimit = 128

	l := list.New(items, list.NewDefaultDelegate(), 0, 10)
	l.Title = "Recent Commands"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return Model{
		Input:      in,
		List:       l,
		Commands:   commands,
		SearchMode: "Exact",
		MatchCount: len(commands),
	}
}
