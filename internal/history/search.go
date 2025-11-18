package history

import (
	"sort"
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

// FuzzyMatch represents a command with its fuzzy match score and positions
type FuzzyMatch struct {
	Command   Command
	Score     int
	Positions []int // Character positions that matched the query
}

// FuzzyMatchResult represents a command with its fuzzy match positions
type FuzzyMatchResult struct {
	Command   Command
	Positions []int // Character positions that matched the query
}

// SearchFuzzy performs fuzzy matching on commands and returns them sorted by relevance
func SearchFuzzy(commands []Command, input string) []Command {
	if strings.TrimSpace(input) == "" {
		return commands
	}

	inputLower := strings.ToLower(input)
	var matches []FuzzyMatch

	for _, cmd := range commands {
		positions := findFuzzyMatchPositions(strings.ToLower(cmd.Text), inputLower)
		if len(positions) == 0 {
			continue
		}
		score := calculateFuzzyScoreWithPositions(cmd.Text, input, inputLower, positions)
		if score > 0 {
			matches = append(matches, FuzzyMatch{
				Command:   cmd,
				Score:     score,
				Positions: positions,
			})
		}
	}

	// Sort by score (higher is better), then by index (lower is better for same score)
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		return matches[i].Command.Index < matches[j].Command.Index
	})

	// Extract commands from sorted matches
	result := make([]Command, len(matches))
	for i, match := range matches {
		result[i] = match.Command
	}

	return result
}

// SearchFuzzyWithPositions performs fuzzy matching and returns commands with their match positions
func SearchFuzzyWithPositions(commands []Command, input string) []FuzzyMatchResult {
	if strings.TrimSpace(input) == "" {
		result := make([]FuzzyMatchResult, len(commands))
		for i, cmd := range commands {
			result[i] = FuzzyMatchResult{
				Command:   cmd,
				Positions: []int{},
			}
		}
		return result
	}

	inputLower := strings.ToLower(input)
	var matches []FuzzyMatch

	for _, cmd := range commands {
		positions := findFuzzyMatchPositions(strings.ToLower(cmd.Text), inputLower)
		if len(positions) == 0 {
			continue
		}
		score := calculateFuzzyScoreWithPositions(cmd.Text, input, inputLower, positions)
		if score > 0 {
			matches = append(matches, FuzzyMatch{
				Command:   cmd,
				Score:     score,
				Positions: positions,
			})
		}
	}

	// Sort by score (higher is better), then by index (lower is better for same score)
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		return matches[i].Command.Index < matches[j].Command.Index
	})

	// Extract commands with positions from sorted matches
	result := make([]FuzzyMatchResult, len(matches))
	for i, match := range matches {
		result[i] = FuzzyMatchResult{
			Command:   match.Command,
			Positions: match.Positions,
		}
	}

	return result
}

// calculateFuzzyScore calculates a relevance score for fuzzy matching
// Returns 0 if the query doesn't match, or a positive score indicating match quality
// This is a wrapper that computes positions internally for backward compatibility
func calculateFuzzyScore(text, query, queryLower string) int {
	textLower := strings.ToLower(text)
	matchPositions := findFuzzyMatchPositions(textLower, queryLower)
	if len(matchPositions) == 0 {
		return 0
	}
	return calculateFuzzyScoreWithPositions(text, query, queryLower, matchPositions)
}

// calculateFuzzyScoreWithPositions calculates a relevance score for fuzzy matching using pre-computed positions
func calculateFuzzyScoreWithPositions(text, query, queryLower string, matchPositions []int) int {
	// Base score: all characters matched
	score := 1000

	// Bonus for consecutive matches (adjacent characters in query are adjacent in text)
	consecutiveBonus := 0
	for i := 1; i < len(matchPositions); i++ {
		if matchPositions[i] == matchPositions[i-1]+1 {
			consecutiveBonus += 50
		}
	}
	score += consecutiveBonus

	// Bonus for matches at the start of the text
	startBonus := 0
	if matchPositions[0] == 0 {
		startBonus = 100
	} else if matchPositions[0] < 5 {
		startBonus = 50 - (matchPositions[0] * 10)
	}
	score += startBonus

	// Bonus for case-sensitive matches
	caseBonus := 0
	for i, pos := range matchPositions {
		if text[pos] == query[i] {
			caseBonus += 5
		}
	}
	score += caseBonus

	// Penalty for spread (distance between first and last match)
	spread := matchPositions[len(matchPositions)-1] - matchPositions[0]
	if spread > len(query) {
		penalty := (spread - len(query)) * 2
		score -= penalty
	}

	// Ensure score is positive
	if score < 1 {
		score = 1
	}

	return score
}

// findFuzzyMatchPositions finds positions where each query character appears in order
// Returns empty slice if not all characters can be matched in order
func findFuzzyMatchPositions(text, query string) []int {
	if len(query) == 0 {
		return []int{}
	}

	positions := make([]int, 0, len(query))
	textIndex := 0

	for _, queryChar := range query {
		found := false
		for textIndex < len(text) {
			if text[textIndex] == byte(queryChar) {
				positions = append(positions, textIndex)
				textIndex++
				found = true
				break
			}
			textIndex++
		}
		if !found {
			return []int{} // Not all characters matched
		}
	}

	return positions
}
