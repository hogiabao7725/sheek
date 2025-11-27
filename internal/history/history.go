package history

import "time"

// CommandSource identifies where a command came from.
type CommandSource string

const (
	CommandSourceUnknown CommandSource = "unknown"
	CommandSourceSheek   CommandSource = "sheek"
)

// CommandContext stores execution metadata for scoring.
type CommandContext struct {
	Directory string
	Repo      string
	Branch    string
	Workspace string
}

// Command represents a shell command with optional metadata.
type Command struct {
	Index        int
	Text         string
	Timestamp    time.Time
	Context      CommandContext
	Source       CommandSource
	ContextBoost int
}

// LoadCommands loads structured sheek history from ~/.sheek_history.
func LoadCommands() ([]Command, error) {
	return LoadSheekHistory()
}
