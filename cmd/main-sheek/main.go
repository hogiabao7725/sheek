package main

import (
	"fmt"
	"log"
	"sheek/internal/history"
	"sheek/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
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

	// Save cursor position before starting TUI (for cleanup)
	fmt.Print("\033[s") // Save cursor position

	// Run inline in current terminal session (like fzf) instead of alt screen
	p := tea.NewProgram(teaModel(model))

	// Run the program
	m, err := p.Run()
	if err != nil {
		log.Fatalf("error running program: %v", err)
	}

	// Clear TUI output and print selected command (like fzf)
	if tm, ok := m.(teaModel); ok {
		selectedModel := tui.Model(tm)
		
		// Restore cursor to saved position (start of TUI)
		fmt.Print("\033[u") // Restore cursor position
		
		// Move cursor up many lines and clear from there to end of screen
		// This ensures we clear all TUI output regardless of how many lines it used
		fmt.Print("\033[50A") // Move up 50 lines (should be enough)
		fmt.Print("\033[0J")  // Clear from cursor to end of screen
		
		if selectedModel.SelectedCommand != "" {
			// Print only the selected command (clean output like fzf)
			fmt.Println(selectedModel.SelectedCommand)
		}
		// If no command selected (ESC/Ctrl+C), nothing is printed (like fzf)
	}
}
