package tui

import (
	"testing"

	"sheek/internal/history"
)

func TestNewModelInitialQueryPrefillsInput(t *testing.T) {
	commands := []history.Command{
		{Index: 0, Text: "git status"},
		{Index: 1, Text: "ls"},
	}

	model := NewModel(commands, "git")

	if model.Input.Value() != "git" {
		t.Fatalf("expected input to be prefilled with 'git', got %q", model.Input.Value())
	}

	if len(model.FilteredCommands) != 1 {
		t.Fatalf("expected filtered commands to contain 1 item, got %d", len(model.FilteredCommands))
	}

	if model.FilteredCommands[0].Text != "git status" {
		t.Fatalf("expected filtered command to be 'git status', got %q", model.FilteredCommands[0].Text)
	}
}
