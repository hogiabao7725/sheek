package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sheek/internal/config"
	"sheek/internal/history"
	"sheek/internal/tui"
	"sheek/internal/tui/styles"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
)

type teaModel tui.Model

func (m teaModel) Init() tea.Cmd { return tui.TickCmd() }
func (m teaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	nm, cmd := tui.Update(msg, tui.Model(m))
	return teaModel(nm), cmd
}
func (m teaModel) View() string { return tui.View(tui.Model(m)) }

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "record":
			if err := handleRecordCommand(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "sheek record: %v\n", err)
				os.Exit(1)
			}
			return
		case "import":
			if err := handleImportCommand(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "sheek import: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	queryFlag := flag.String("query", "", "prefill the search input with a query")
	flag.Parse()

	initialQuery := *queryFlag
	if initialQuery == "" {
		initialQuery = os.Getenv("SHEEK_INITIAL_QUERY")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize styles with config colors
	styles.InitializeStyles(cfg)

	cmds, err := history.LoadCommands()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load history: %v\n", err)
		os.Exit(1)
	}

	ctx := history.CurrentContext()
	cmds = history.ApplyContextBoost(cmds, ctx)

	model := tui.NewModel(cmds, cfg, initialQuery)

	// Ensure TERM is set for color support (important when running from keybind)
	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
		os.Setenv("TERM", term)
	}

	// Ensure NO_COLOR is not set (it would disable colors)
	if os.Getenv("NO_COLOR") != "" {
		os.Unsetenv("NO_COLOR")
	}

	// Force color mode - critical when running from shell widget where auto-detection may fail
	// Check TERM first as it's more reliable than TTY detection in shell widgets
	isTTY := isatty.IsTerminal(os.Stderr.Fd())

	// Determine color profile based on TERM and TTY status
	// When running from shell widget, stderr may not be a TTY but terminal still supports colors
	if term == "xterm-256color" || term == "screen-256color" || term == "tmux-256color" ||
		term == "xterm-kitty" || term == "wezterm" || isTTY {
		// Force true color - most modern terminals support it
		os.Setenv("COLORTERM", "truecolor")
		lipgloss.SetColorProfile(termenv.TrueColor)
	} else if term != "" && (term == "xterm" || term == "screen" || term == "tmux") {
		// Fallback to 256 colors for basic color terminals
		os.Setenv("COLORTERM", "256color")
		lipgloss.SetColorProfile(termenv.ANSI256)
	} else {
		// Last resort: try 256 colors anyway
		lipgloss.SetColorProfile(termenv.ANSI256)
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
		// If no command selected (ESC/Ctrl+C), exit with non-zero code
		os.Exit(1)
	}
}

func handleRecordCommand(args []string) error {
	fs := flag.NewFlagSet("record", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cmdText := fs.String("cmd", "", "command text to record")
	cwd := fs.String("cwd", "", "command working directory")
	repo := fs.String("repo", "", "repository override")
	branch := fs.String("branch", "", "branch override")
	workspace := fs.String("workspace", "", "workspace override")
	ts := fs.Int64("ts", 0, "unix timestamp (seconds)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*cmdText) == "" {
		return fmt.Errorf("--cmd is required")
	}

	payload := history.RecordPayload{
		Command:   *cmdText,
		Directory: *cwd,
		Repo:      *repo,
		Branch:    *branch,
		Workspace: *workspace,
	}

	if *ts > 0 {
		payload.Timestamp = time.Unix(*ts, 0)
	}

	return history.RecordCommand(payload)
}

func handleImportCommand(args []string) error {
	fs := flag.NewFlagSet("import", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	shell := fs.String("shell", "", "shell to import (zsh|bash|fish)")
	source := fs.String("source", "", "path to custom history file")
	limit := fs.Int("limit", 0, "maximum commands to import (0 = all)")
	appendMode := fs.Bool("append", true, "append instead of replacing history file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*shell) == "" {
		return fmt.Errorf("--shell is required")
	}

	count, err := history.ImportHistory(history.ImportOptions{
		Shell:      *shell,
		SourcePath: *source,
		Limit:      *limit,
		Append:     *appendMode,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Imported %d commands into sheek history\n", count)
	return nil
}
