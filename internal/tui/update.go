package tui

import (
	"time"

	"sheek/internal/history"
	"sheek/internal/tui/components"

	tea "github.com/charmbracelet/bubbletea"
)

const relativeTimeRefreshInterval = time.Second

type tickMsg struct{}

// TickCmd returns a command that schedules periodic ticks for refreshing relative timestamps.
func TickCmd() tea.Cmd {
	return newTickCmd()
}

func newTickCmd() tea.Cmd {
	return tea.Tick(relativeTimeRefreshInterval, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Update handles application updates based on messages
func Update(msg tea.Msg, model Model) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model = handleWindowResize(model, msg)
	case tea.KeyMsg:
		if shouldQuit(msg.String()) {
			return model, tea.Quit
		}
		if msg.String() == "tab" {
			model.SearchMode = model.SearchMode.Toggle()
		}
		if msg.String() == "enter" {
			return handleEnterKey(model)
		}
		if isNavigationKey(msg.String()) {
			var listCmd tea.Cmd
			model.List, listCmd = model.List.Update(msg)
			cmds = append(cmds, listCmd)
		}
	case tickMsg:
		cmds = append(cmds, newTickCmd())
	}

	// Update input
	var cmd tea.Cmd
	model.Input, cmd = model.Input.Update(msg)
	cmds = append(cmds, cmd)

	// Update search results
	model = updateSearchResults(model)

	return model, tea.Batch(cmds...)
}

// handleWindowResize updates model dimensions when window is resized
func handleWindowResize(model Model, msg tea.WindowSizeMsg) Model {
	model.Width = msg.Width
	model.Height = msg.Height
	model.List.SetSize(msg.Width-4, msg.Height-10)
	return model
}

// shouldQuit returns true if the key should quit the application
func shouldQuit(key string) bool {
	return key == "ctrl+c" || key == "esc"
}

// isNavigationKey returns true if the key is used for navigation
func isNavigationKey(key string) bool {
	return key == "up" || key == "down"
}

// handleEnterKey handles the enter key press
func handleEnterKey(model Model) (Model, tea.Cmd) {
	// Save the selected command before quitting (like fzf output)
	// This command will be printed to stdout in main.go
	if len(model.FilteredCommands) > 0 {
		selectedIndex := model.List.Index()
		if selectedIndex >= 0 && selectedIndex < len(model.FilteredCommands) {
			model.SelectedCommand = model.FilteredCommands[selectedIndex].Text
		}
	}
	// Quit the program - main.go will handle printing the selected command
	return model, tea.Quit
}

// updateSearchResults performs search filtering and updates the model
func updateSearchResults(model Model) Model {
	inputValue := model.Input.Value()
	var filtered []history.Command

	switch model.SearchMode {
	case SearchModeExact:
		filtered = history.SearchExact(model.Commands, inputValue)
		model.FuzzyPositions = nil // Clear fuzzy positions for exact search
	case SearchModeFuzzy:
		fuzzyResults := history.SearchFuzzyWithPositions(model.Commands, inputValue)
		// Build map from command index to match positions
		model.FuzzyPositions = make(map[int][]int, len(fuzzyResults))
		filtered = make([]history.Command, len(fuzzyResults))
		for i, result := range fuzzyResults {
			filtered[i] = result.Command
			model.FuzzyPositions[result.Command.Index] = result.Positions
		}
	default:
		filtered = history.SearchExact(model.Commands, inputValue)
		model.FuzzyPositions = nil
	}

	// Check if search results changed
	previousCount := len(model.FilteredCommands)
	searchResultsChanged := previousCount != len(filtered)

	model.FilteredCommands = filtered
	model.List.SetItems(components.CommandsToListItems(filtered))

	// Adjust list index if search results changed
	if searchResultsChanged {
		model = adjustListIndex(model, len(filtered))
	}

	return model
}

// adjustListIndex ensures the list index is valid after search results change
func adjustListIndex(model Model, resultCount int) Model {
	currentIndex := model.List.Index()

	if resultCount == 0 {
		model.List.Select(0)
	} else if currentIndex >= resultCount {
		model.List.Select(0)
	} else if currentIndex < resultCount {
		model.List.Select(currentIndex)
	}

	return model
}
