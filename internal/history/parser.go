package history

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var zshHistoryPrefix = regexp.MustCompile(`^: [0-9]+:[0-9]+;`)

// Save current command if not empty, then reset builder.
func flushCurrent(commands *[]Command, builder *strings.Builder, index *int, timestamp time.Time) {
	text := strings.TrimSpace(builder.String())
	if text != "" {
		*commands = append(*commands, Command{
			Index:     *index,
			Text:      text,
			Timestamp: timestamp,
		})
		*index++
	}
	builder.Reset()
}

// ParseZshHistory converts raw Zsh history lines into commands.
func ParseZshHistory(rawLines []string) []Command {
	var (
		commands         []Command
		current          strings.Builder
		index            = 1
		currentTimestamp time.Time
	)

	for _, line := range rawLines {
		if zshHistoryPrefix.MatchString(line) {
			if current.Len() > 0 {
				flushCurrent(&commands, &current, &index, currentTimestamp)
			}

			ts, cmd := extractMetadataAndCommand(line)
			currentTimestamp = ts
			if cmd != "" {
				current.WriteString(cmd)
			}
		} else {
			if current.Len() > 0 {
				current.WriteString("\n")
			}
			current.WriteString(line)
		}
	}

	if current.Len() > 0 {
		flushCurrent(&commands, &current, &index, currentTimestamp)
	}
	return commands
}

func extractMetadataAndCommand(line string) (time.Time, string) {
	parts := strings.SplitN(line, ";", 2)
	if len(parts) != 2 {
		return time.Time{}, ""
	}

	meta := strings.TrimPrefix(parts[0], ": ")
	metaParts := strings.Split(meta, ":")
	if len(metaParts) == 0 {
		return time.Time{}, parts[1]
	}

	epoch, err := strconv.ParseInt(metaParts[0], 10, 64)
	if err != nil {
		return time.Time{}, parts[1]
	}

	return time.Unix(epoch, 0), parts[1]
}
