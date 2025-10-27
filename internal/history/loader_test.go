package history

import (
	"testing"
)

func TestLoadAndParseZshHistory(t *testing.T) {
	lines, err := LoadAndParseZshHistory()
	if err != nil {
		t.Fatalf("Failed to load history: %v", err)
	}
	if len(lines) == 0 {
		t.Errorf("Expected some commands, got zero")
	}
	t.Logf("Loaded %d commands from history", len(lines))
}
