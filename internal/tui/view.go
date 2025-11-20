package tui

import (
	"strings"

	"sheek/internal/tui/components"
	"sheek/internal/tui/styles"
)

// View renders the application UI
func View(model Model) string {
	var b strings.Builder

	inputContent := model.Input.View()
	if model.Input.Value() == "" && model.Placeholder != "" {
		placeholder := styles.SearchPlaceholderStyle.Render(model.Placeholder)
		inputContent = inputContent + placeholder
	}

	searchBar := components.RenderSearchComponent("> ", inputContent, model.SearchMode.String(), model.Width)
	b.WriteString(searchBar)

	listView := components.RenderListComponent(model.FilteredCommands, model.FuzzyPositions, model.List.Index(), model.Width, model.Height, model.Input.Value(), components.SearchMode(model.SearchMode))
	b.WriteString(listView)
	b.WriteString("\n")

	return b.String()
}
