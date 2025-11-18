package main

import (
	"fmt"
	"os"
	"sheek/internal/history"
	"sheek/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
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
		fmt.Fprintf(os.Stderr, "failed to load history: %v\n", err)
		os.Exit(1)
	}

	model := tui.NewModel(cmds)

	// Ensure TERM is set for color support (important when running from keybind)
	if os.Getenv("TERM") == "" {
		os.Setenv("TERM", "xterm-256color")
	}

	// Force color mode if stderr is a TTY
	// This ensures colors work even when called from shell keybind
	if isatty.IsTerminal(os.Stderr.Fd()) {
		os.Setenv("COLORTERM", "truecolor")
	}

	// Save cursor position before starting TUI (for cleanup)
	fmt.Fprint(os.Stderr, "\033[s") // Save cursor position

	// Run inline in current terminal session (like fzf) instead of alt screen
	// Redirect bubbletea output to stderr so stdout is clean for command output
	p := tea.NewProgram(teaModel(model), tea.WithOutput(os.Stderr))

	// Run the program
	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running program: %v\n", err)
		os.Exit(1)
	}

		// Clear TUI output and print selected command (like fzf)
		if tm, ok := m.(teaModel); ok {
			selectedModel := tui.Model(tm)

			// Restore cursor to saved position (start of TUI)
			fmt.Fprint(os.Stderr, "\033[u") // Restore cursor position

			// Move cursor up many lines and clear from there to end of screen
			// This ensures we clear all TUI output regardless of how many lines it used
			fmt.Fprint(os.Stderr, "\033[50A") // Move up 50 lines (should be enough)
			fmt.Fprint(os.Stderr, "\033[0J")  // Clear from cursor to end of screen

			if selectedModel.SelectedCommand != "" {
				// Print command to stderr so it's visible in terminal before prompt
				// This allows user to see the selected command (like fzf behavior)
				fmt.Fprintln(os.Stderr, selectedModel.SelectedCommand)
				
				// Also print to stdout for shell integration (keybinds scripts)
				// This allows shell to capture and use the command
				fmt.Fprintln(os.Stdout, selectedModel.SelectedCommand)
				os.Exit(0)
			}
			// If no command selected (ESC/Ctrl+C), exit with non-zero code (like fzf)
			os.Exit(1)
		}
}
