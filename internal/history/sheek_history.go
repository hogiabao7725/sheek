package history

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const sheekHistoryFilename = ".sheek_history"

// ErrSheekHistoryMissing indicates there is no structured history yet.
var ErrSheekHistoryMissing = errors.New("sheek history not found")

type historyRecord struct {
	Command   string `json:"cmd"`
	Timestamp int64  `json:"ts"`
	Directory string `json:"cwd"`
	Repo      string `json:"repo,omitempty"`
	Branch    string `json:"branch,omitempty"`
	Workspace string `json:"workspace,omitempty"`
}

// RecordPayload is used by shell hooks to append new history entries.
type RecordPayload struct {
	Command   string
	Directory string
	Repo      string
	Branch    string
	Workspace string
	Timestamp time.Time
}

func sheekHistoryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".config", "sheek")
	return filepath.Join(configDir, sheekHistoryFilename), nil
}

// LoadSheekHistory reads ~/.sheek_history JSONL entries.
func LoadSheekHistory() ([]Command, error) {
	path, err := sheekHistoryPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrSheekHistoryMissing
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var (
		commands []Command
		index    = 1
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var record historyRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			continue // Skip malformed entries silently
		}

		ts := record.Timestamp
		if ts < 1 {
			ts = 1
		}

		commands = append(commands, Command{
			Index:     index,
			Text:      record.Command,
			Timestamp: time.Unix(ts, 0),
			Context: CommandContext{
				Directory: record.Directory,
				Repo:      record.Repo,
				Branch:    record.Branch,
				Workspace: record.Workspace,
			},
			Source: CommandSourceSheek,
		})
		index++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return commands, nil
}

// RecordCommand appends a new entry to ~/.sheek_history.
func RecordCommand(payload RecordPayload) error {
	command := strings.TrimSpace(payload.Command)
	if command == "" {
		return fmt.Errorf("cannot record empty command")
	}

	directory := payload.Directory
	if directory == "" {
		var err error
		directory, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	absDir, err := filepath.Abs(directory)
	if err != nil {
		absDir = directory
	}
	absDir = filepath.Clean(absDir)

	ts := payload.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}
	if ts.Unix() < 1 {
		ts = time.Unix(1, 0)
	}

	ctx := ResolveContext(absDir)
	if payload.Repo != "" {
		ctx.Repo = payload.Repo
	}
	if payload.Branch != "" {
		ctx.Branch = payload.Branch
	}
	if payload.Workspace != "" {
		ctx.Workspace = payload.Workspace
	}

	record := historyRecord{
		Command:   command,
		Timestamp: ts.Unix(),
		Directory: absDir,
		Repo:      ctx.Repo,
		Branch:    ctx.Branch,
		Workspace: ctx.Workspace,
	}

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	path, err := sheekHistoryPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(append(data, '\n')); err != nil {
		return err
	}

	return nil
}

