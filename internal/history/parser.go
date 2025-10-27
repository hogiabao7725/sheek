package history

import (
	"regexp"
	"strings"
)

var zshHistoryPrefix = regexp.MustCompile(`^: [0-9]+:[0-9]+;`)

// Save current command if not empty, then reset builder.
func flushCurrent(commands *[]Command, builder *strings.Builder, index *int) {
	text := strings.TrimSpace(builder.String())
	if text != "" {
		*commands = append(*commands, Command{Index: *index, Text: text})
		*index++
	}
	builder.Reset()
}

// ParseZshHistory converts raw Zsh history lines into commands.
func ParseZshHistory(rawLines []string) []Command {
	var commands []Command
	var current strings.Builder
	index := 1

	for _, line := range rawLines {
		if zshHistoryPrefix.MatchString(line) {
			if current.Len() > 0 {
				flushCurrent(&commands, &current, &index)
			}
			parts := strings.SplitN(line, ";", 2)
			if len(parts) == 2 {
				current.WriteString(parts[1])
			}
		} else {
			if current.Len() > 0 {
				current.WriteString("\n")
			}
			current.WriteString(line)
		}
	}

	if current.Len() > 0 {
		flushCurrent(&commands, &current, &index)
	}
	return commands
}
