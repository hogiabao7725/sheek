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

	// Pass config values to search component
	searchBar := components.RenderSearchComponent(
		"> ",
		inputContent,
		model.SearchMode.String(),
		model.Width,
		model.Config.Margin,
	)
	b.WriteString(searchBar)
	b.WriteString("\n")

	// Pass config values to list component
	listView := components.RenderListComponent(
		model.FilteredCommands,
		model.FuzzyPositions,
		model.List.Index(),
		model.Width,
		model.Height,
		model.Input.Value(),
		components.SearchMode(model.SearchMode),
		model.Config.MaxItems,
		model.Config.Height,
		model.Config.Margin,
	)
	b.WriteString(listView)
	b.WriteString("\n")

	return b.String()
}
