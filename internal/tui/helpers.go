package tui

import (
	"fmt"
	"sheek/internal/history"

	"github.com/charmbracelet/bubbles/list"
)

type commandItem struct {
	command history.Command
}

func (c commandItem) Title() string {
	return fmt.Sprintf("%d", c.command.Index)
}

func (c commandItem) Description() string {
	return c.command.Text
}

func (c commandItem) FilterValue() string {
	return c.command.Text
}

func CommandsToListItems(commands []history.Command) []list.Item {
	items := make([]list.Item, len(commands))
	for i, cmd := range commands {
		items[i] = commandItem{command: cmd}
	}
	return items
}
