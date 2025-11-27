package history

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ImportOptions configures shell history import.
type ImportOptions struct {
	Shell      string
	SourcePath string
	Limit      int
	Append     bool
}

// ImportHistory imports shell history into the sheek history file.
func ImportHistory(opts ImportOptions) (int, error) {
	shell := strings.ToLower(strings.TrimSpace(opts.Shell))
	if shell == "" {
		return 0, fmt.Errorf("import shell is required")
	}

	path, err := resolveHistoryPath(shell, opts.SourcePath)
	if err != nil {
		return 0, err
	}

	var records []HistoryRecord
	switch shell {
	case "zsh":
		records, err = importZshHistory(path, opts.Limit)
	case "bash":
		records, err = importBashHistory(path, opts.Limit)
	case "fish":
		records, err = importFishHistory(path, opts.Limit)
	default:
		return 0, fmt.Errorf("unsupported shell %q", shell)
	}
	if err != nil {
		return 0, err
	}

	if len(records) == 0 {
		return 0, errors.New("no commands were imported")
	}

	if opts.Append {
		if err := AppendHistoryRecords(records); err != nil {
			return 0, err
		}
	} else {
		if err := ReplaceHistoryRecords(records); err != nil {
			return 0, err
		}
	}

	return len(records), nil
}

func resolveHistoryPath(shell, customPath string) (string, error) {
	if customPath != "" {
		return customPath, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch shell {
	case "zsh":
		return filepath.Join(home, ".zsh_history"), nil
	case "bash":
		return filepath.Join(home, ".bash_history"), nil
	case "fish":
		return filepath.Join(home, ".local", "share", "fish", "fish_history"), nil
	default:
		return "", fmt.Errorf("no default history path for shell %q", shell)
	}
}

func importZshHistory(path string, limit int) ([]HistoryRecord, error) {
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}

	records := parseZshHistory(lines)
	return trimRecords(records, limit), nil
}

var zshHistoryPrefix = regexp.MustCompile(`^: ([0-9]+):[0-9]+;`)

func parseZshHistory(lines []string) []HistoryRecord {
	var (
		current strings.Builder
		records []HistoryRecord
		ts      int64
	)

	flush := func() {
		text := strings.TrimSpace(current.String())
		if text == "" {
			return
		}
		records = append(records, HistoryRecord{
			Command:   text,
			Timestamp: ensureTimestamp(ts),
		})
		current.Reset()
	}

	for _, line := range lines {
		if matches := zshHistoryPrefix.FindStringSubmatch(line); len(matches) == 2 {
			if current.Len() > 0 {
				flush()
			}
			parsed, _ := strconv.ParseInt(matches[1], 10, 64)
			ts = parsed
			sep := strings.IndexRune(line, ';')
			if sep >= 0 && sep+1 < len(line) {
				current.WriteString(line[sep+1:])
			}
		} else {
			if current.Len() > 0 {
				current.WriteRune('\n')
			}
			current.WriteString(line)
		}
	}

	if current.Len() > 0 {
		flush()
	}

	return records
}

func importBashHistory(path string, limit int) ([]HistoryRecord, error) {
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}

	records := make([]HistoryRecord, 0, len(lines))
	var baseTS = time.Now().Unix()
	for i := len(lines) - 1; i >= 0; i-- {
		text := strings.TrimSpace(lines[i])
		if text == "" {
			continue
		}
		records = append(records, HistoryRecord{
			Command:   text,
			Timestamp: ensureTimestamp(baseTS),
		})
		baseTS--
	}

	reverseRecords(records)
	return trimRecords(records, limit), nil
}

func importFishHistory(path string, limit int) ([]HistoryRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []HistoryRecord
	var current HistoryRecord
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "- cmd:") {
			if current.Command != "" {
				records = append(records, current)
				current = HistoryRecord{}
			}
			current.Command = strings.TrimSpace(strings.TrimPrefix(line, "- cmd:"))
		} else if strings.HasPrefix(line, "when:") {
			// old format
			ts, _ := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(line, "when:")), 10, 64)
			current.Timestamp = ensureTimestamp(ts)
		} else if strings.HasPrefix(line, "time:") {
			ts, _ := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(line, "time:")), 10, 64)
			current.Timestamp = ensureTimestamp(ts)
		}
	}
	if current.Command != "" {
		records = append(records, current)
	}

	for i := range records {
		if records[i].Timestamp == 0 {
			records[i].Timestamp = ensureTimestamp(time.Now().Unix())
		}
	}

	return trimRecords(records, limit), nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func trimRecords(records []HistoryRecord, limit int) []HistoryRecord {
	if limit <= 0 || limit >= len(records) {
		return records
	}
	return records[len(records)-limit:]
}

func reverseRecords(records []HistoryRecord) {
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}
}

func ensureTimestamp(value int64) int64 {
	if value < 1 {
		return 1
	}
	return value
}
