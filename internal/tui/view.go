package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func View(model Model) string {
	var b strings.Builder

	searchBar := RenderSearchComponent("> ", model.Input.Value(), model.SearchMode, model.Width)
	b.WriteString(searchBar)
	b.WriteString("\n")

	listView := RenderListComponent(model.Commands, model.List.Index(), model.Width, model.Height)
	b.WriteString(listView)
	b.WriteString("\n")

	helpText := "tab: toggle mode │ enter: select │ esc: quit"
	helpStyle := lipgloss.NewStyle().Foreground(mutedColor).MarginLeft(2)
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}
