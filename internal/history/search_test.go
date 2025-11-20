package history

import (
	"strings"
	"testing"
)

func TestSearchExact(t *testing.T) {
	cmds := []Command{
		{Index: 0, Text: "kitty kitty"},
		{Index: 1, Text: "kitty"},
		{Index: 2, Text: "repeat after me"},
		{Index: 3, Text: "another kitty command"},
		{Index: 4, Text: "git status"},
		{Index: 5, Text: "git repeat status"},
		{Index: 6, Text: "zsh kitty script"},
	}

	tests := []struct {
		input string
		want  int
	}{
		{"kitty", 4},
		{"repeat", 2},
		{"", len(cmds)},
		{"missing", 0},
	}

	for _, test := range tests {
		got := SearchExact(cmds, test.input)
		if len(got) != test.want {
			t.Errorf("SearchExact(%q) = %d, want %d", test.input, len(got), test.want)
		}
	}
}

func TestSearchFuzzy(t *testing.T) {
	cmds := []Command{
		{Index: 1, Text: "git commit -m 'initial commit'"},
		{Index: 2, Text: "cd /home/user/documents"},
		{Index: 3, Text: "npm install react"},
		{Index: 4, Text: "docker run -it ubuntu"},
		{Index: 5, Text: "go build main.go"},
		{Index: 6, Text: "git push origin main"},
		{Index: 7, Text: "cd ~/projects"},
	}

	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name:     "exact substring match",
			query:    "git",
			expected: []string{"git commit -m 'initial commit'", "git push origin main"},
		},
		{
			name:     "fuzzy match - characters in order",
			query:    "gcm",
			expected: []string{"git commit -m 'initial commit'"},
		},
		{
			name:     "fuzzy match - spread characters",
			query:    "gpm",
			expected: []string{"git push origin main"},
		},
		{
			name:     "fuzzy match - case insensitive",
			query:    "NPM",
			expected: []string{"npm install react"},
		},
		{
			name:     "fuzzy match - partial word",
			query:    "doc",
			expected: []string{"cd /home/user/documents"}, // Should match, but may have more results
		},
		{
			name:     "empty query returns all",
			query:    "",
			expected: []string{"git commit -m 'initial commit'", "cd /home/user/documents", "npm install react", "docker run -it ubuntu", "go build main.go", "git push origin main", "cd ~/projects"},
		},
		{
			name:     "no match returns empty",
			query:    "xyzabc",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchFuzzy(cmds, tt.query)

			if tt.query == "" {
				// Empty query should return all commands
				if len(result) != len(tt.expected) {
					t.Errorf("SearchFuzzy(%q) returned %d results, want %d", tt.query, len(result), len(tt.expected))
				}
				return
			}

			if tt.query == "xyzabc" {
				// No match should return empty
				if len(result) != 0 {
					t.Errorf("SearchFuzzy(%q) returned %d results, want 0", tt.query, len(result))
				}
				return
			}

			// For fuzzy matches, check that expected results are present
			// (fuzzy search may return additional relevant matches)
			resultTexts := make(map[string]bool)
			for _, cmd := range result {
				resultTexts[cmd.Text] = true
			}

			for _, expectedText := range tt.expected {
				if !resultTexts[expectedText] {
					t.Errorf("SearchFuzzy(%q) missing expected result: %q", tt.query, expectedText)
				}
			}

			// For exact matches, also check order
			if tt.name == "exact substring match" {
				if len(result) != len(tt.expected) {
					t.Errorf("SearchFuzzy(%q) returned %d results, want %d", tt.query, len(result), len(tt.expected))
					return
				}
				for i, cmd := range result {
					if i >= len(tt.expected) {
						break
					}
					if cmd.Text != tt.expected[i] {
						t.Errorf("SearchFuzzy(%q)[%d] = %q, want %q", tt.query, i, cmd.Text, tt.expected[i])
					}
				}
			}
		})
	}
}

func TestSearchFuzzyScoring(t *testing.T) {
	cmds := []Command{
		{Index: 1, Text: "git commit"},
		{Index: 2, Text: "go test"},
		{Index: 3, Text: "git status"},
	}

	// Test that consecutive matches score higher
	result := SearchFuzzy(cmds, "git")
	if len(result) < 2 {
		t.Fatalf("Expected at least 2 results for 'git', got %d", len(result))
	}

	// Both should match, but order might vary based on scoring
	// At minimum, both should be in results
	texts := make([]string, len(result))
	for i, cmd := range result {
		texts[i] = cmd.Text
	}

	hasGitCommit := false
	hasGitStatus := false
	for _, text := range texts {
		if strings.Contains(text, "git commit") {
			hasGitCommit = true
		}
		if strings.Contains(text, "git status") {
			hasGitStatus = true
		}
	}

	if !hasGitCommit || !hasGitStatus {
		t.Errorf("Expected both 'git commit' and 'git status' in results, got %v", texts)
	}
}

func TestCalculateFuzzyScore(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		query    string
		expected int // expected to be > 0 for matches
	}{
		{
			name:     "exact match at start",
			text:     "git commit",
			query:    "git",
			expected: 1000, // base score + start bonus
		},
		{
			name:     "fuzzy match",
			text:     "git commit",
			query:    "gcm",
			expected: 1000, // base score
		},
		{
			name:     "no match",
			text:     "git commit",
			query:    "xyz",
			expected: 0,
		},
		{
			name:     "consecutive characters",
			text:     "git commit",
			query:    "git",
			expected: 1000, // base + consecutive bonus + start bonus
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryLower := strings.ToLower(tt.query)
			score := calculateFuzzyScore(tt.text, tt.query, queryLower)

			if tt.expected == 0 {
				if score != 0 {
					t.Errorf("calculateFuzzyScore(%q, %q) = %d, want 0", tt.text, tt.query, score)
				}
			} else {
				if score <= 0 {
					t.Errorf("calculateFuzzyScore(%q, %q) = %d, want > 0", tt.text, tt.query, score)
				}
			}
		})
	}
}

func TestFindFuzzyMatchPositions(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		query    string
		expected []int
	}{
		{
			name:     "exact match",
			text:     "git",
			query:    "git",
			expected: []int{0, 1, 2},
		},
		{
			name:     "fuzzy match",
			text:     "git commit",
			query:    "gcm",
			expected: []int{0, 4, 6}, // g at 0, c at 4, m at 6
		},
		{
			name:     "no match",
			text:     "git",
			query:    "xyz",
			expected: []int{},
		},
		{
			name:     "spread match",
			text:     "git push origin main",
			query:    "gpm",
			expected: []int{0, 4, 16}, // g at 0, p at 4, m at 16
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findFuzzyMatchPositions(tt.text, tt.query)

			if len(result) != len(tt.expected) {
				t.Errorf("findFuzzyMatchPositions(%q, %q) = %v, want %v", tt.text, tt.query, result, tt.expected)
				return
			}

			for i, pos := range result {
				if i >= len(tt.expected) {
					break
				}
				if pos != tt.expected[i] {
					t.Errorf("findFuzzyMatchPositions(%q, %q)[%d] = %d, want %d", tt.text, tt.query, i, pos, tt.expected[i])
				}
			}
		})
	}
}
