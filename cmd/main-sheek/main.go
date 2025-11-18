package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"sheek/internal/history"
	"sheek/internal/tui"
)

type teaModel tui.Model

func (m teaModel) Init() tea.Cmd { return nil }
func (m teaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	nm, cmd := tui.Update(msg, tui.Model(m))
	return teaModel(nm), cmd
}
func (m teaModel) View() string { return tui.View(tui.Model(m)) }

func main() {
	cmds, err := history.LoadAndParseZshHistory()
	if err != nil {
		log.Fatalf("failed to load history: %v", err)
	}

	model := tui.NewModel(cmds)

	// Use AltScreen to draw full screen without overwriting old prompt
	p := tea.NewProgram(teaModel(model), tea.WithAltScreen())

	// Run the program and clean up terminal after quit
	if _, err := p.Run(); err != nil {
		log.Fatalf("error running program: %v", err)
	}

	// Clear screen after program ends (fzf style)
	print("\033[H\033[2J")
}
