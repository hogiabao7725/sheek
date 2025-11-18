package components

import (
	"fmt"

	"sheek/internal/history"

	"github.com/charmbracelet/bubbles/list"
)

// commandItem wraps a history.Command to implement the list.Item interface
type commandItem struct {
	command history.Command
}

// Title returns the command index as a string
func (c commandItem) Title() string {
	return fmt.Sprintf("%d", c.command.Index)
}

// Description returns the command text
func (c commandItem) Description() string {
	return c.command.Text
}

// FilterValue returns the command text for filtering
func (c commandItem) FilterValue() string {
	return c.command.Text
}

// CommandsToListItems converts a slice of commands to list items
func CommandsToListItems(commands []history.Command) []list.Item {
	items := make([]list.Item, len(commands))
	for i, cmd := range commands {
		items[i] = commandItem{command: cmd}
	}
	return items
}

