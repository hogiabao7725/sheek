package tui

import (
	"sheek/internal/history"

	tea "github.com/charmbracelet/bubbletea"
)

func Update(msg tea.Msg, model Model) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.Width = msg.Width
		model.Height = msg.Height
		model.List.SetSize(msg.Width-4, msg.Height-10)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return model, tea.Quit
		case "tab":
			if model.SearchMode == "Exact" {
				model.SearchMode = "Fuzzy"
			} else {
				model.SearchMode = "Exact"
			}
		case "enter":
			selected := model.List.SelectedItem()
			if selected != nil {
				// TODO: xử lý khi người dùng chọn lệnh
			}
			return model, tea.Quit
		}
	}

	var cmd tea.Cmd
	model.Input, cmd = model.Input.Update(msg)
	cmds = append(cmds, cmd)

	inputValue := model.Input.Value()
	var filtered []history.Command
	if model.SearchMode == "Exact" {
		filtered = history.SearchExact(model.Commands, inputValue)
	} else {
		filtered = history.SearchExact(model.Commands, inputValue) // placeholder fuzzy
	}

	model.MatchCount = len(filtered)
	model.List.SetItems(CommandsToListItems(filtered))

	return model, tea.Batch(cmds...)
}
