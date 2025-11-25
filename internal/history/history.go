package history

import "time"

type Command struct {
	Index     int
	Text      string
	Timestamp time.Time
}

func LoadAndParseZshHistory() ([]Command, error) {
	rawLines, err := LoadZshHistory()
	if err != nil {
		return nil, err
	}
	cmds := ParseZshHistory(rawLines)
	return cmds, nil
}
