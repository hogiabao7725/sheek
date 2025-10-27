package history

import (
	"strings"
)

func SearchExact(commands []Command, input string) []Command {
	if strings.TrimSpace(input) == "" {
		return commands
	}

	inputLower := strings.ToLower(input)
	var filtered []Command
	for _, cmd := range commands {
		if strings.Contains(strings.ToLower(cmd.Text), inputLower) {
			filtered = append(filtered, cmd)
		}
	}
	return filtered
}
