package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func View(model Model) string {
	var b strings.Builder

	// Render search component
	searchBar := RenderSearchComponent(
		"> ",
		model.Input.Value(),
		model.SearchMode,
		model.Width,
	)
	b.WriteString(searchBar)
	b.WriteString("\n")

	// Match count
	matchInfo := fmt.Sprintf("%d matches", model.MatchCount)
	infoStyle := lipgloss.NewStyle().
		Foreground(MutedColor).
		MarginLeft(2)
	b.WriteString(infoStyle.Render(matchInfo))
	b.WriteString("\n")

	// Command list
	b.WriteString(model.List.View())
	b.WriteString("\n")

	// Help text
	helpText := "tab: toggle mode │ enter: select │ esc: quit"
	helpStyle := lipgloss.NewStyle().
		Foreground(MutedColor).
		MarginLeft(2)
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}
