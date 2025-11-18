package tui

import (
	"strings"

	"sheek/internal/tui/components"
)

// View renders the application UI
func View(model Model) string {
	var b strings.Builder

	searchBar := components.RenderSearchComponent("> ", model.Input.Value(), model.SearchMode.String(), model.Width)
	b.WriteString(searchBar)

	listView := components.RenderListComponent(model.FilteredCommands, model.FuzzyPositions, model.List.Index(), model.Width, model.Height, model.Input.Value(), components.SearchMode(model.SearchMode))
	b.WriteString(listView)
	b.WriteString("\n")

	return b.String()
}
