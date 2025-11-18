package tui

import (
	"strings"

	"sheek/internal/tui/components"
	"sheek/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// View renders the application UI
func View(model Model) string {
	var b strings.Builder

	searchBar := components.RenderSearchComponent("> ", model.Input.Value(), model.SearchMode.String(), model.Width)
	b.WriteString(searchBar)
	b.WriteString("\n")

	listView := components.RenderListComponent(model.FilteredCommands, model.List.Index(), model.Width, model.Height, model.Input.Value())
	b.WriteString(listView)
	b.WriteString("\n")

	helpText := "tab: toggle mode │ enter: select │ esc: quit"
	helpStyle := lipgloss.NewStyle().Foreground(styles.MutedColor).MarginLeft(2)
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}
