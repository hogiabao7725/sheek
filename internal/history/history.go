package history

type Command struct {
	Index int
	Text  string
}

func LoadAndParseZshHistory() ([]Command, error) {
	rawLines, err := LoadZshHistory()
	if err != nil {
		return nil, err
	}
	cmds := ParseZshHistory(rawLines)
	return cmds, nil
}
