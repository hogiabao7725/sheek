package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"sheek/internal/history"
)

// Colors
var (
	PrimaryColor   = lipgloss.Color("#7D56F4")
	SecondaryColor = lipgloss.Color("#04B575")
	MutedColor     = lipgloss.Color("#626262")
	BorderColor    = lipgloss.Color("#FFFFFF")
	TextColor      = lipgloss.Color("#FFFFFF")
	BackgroundColor = lipgloss.Color("#1A1A1A")
	SelectedColor  = lipgloss.Color("#2D2D2D")
	AccentColor    = lipgloss.Color("#FFD700")
)

// Search Component Styles
var (
	// PromptStyle Prompt symbol (>)
	PromptStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	// Left panel - search input box
	SearchInputBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(BorderColor).
				Foreground(TextColor).
				PaddingLeft(1).
				PaddingRight(1).
				Height(1)

	// Right panel - mode badge
	ModeBadgeStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor).
			Foreground(SecondaryColor).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1).
			Height(1).
			Align(lipgloss.Center)

	// Container for search bar
	SearchContainerStyle = lipgloss.NewStyle().
				MarginLeft(1).
				MarginRight(1).
				MarginTop(1).
				MarginBottom(1)
)

// List Component Styles
var (
	// List container style
	ListContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(BorderColor).
				Background(BackgroundColor).
				MarginLeft(1).
				MarginRight(1).
				MarginTop(1).
				Padding(1)

	// Individual list item styles
	ListItemStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			PaddingLeft(1).
			PaddingRight(1).
			PaddingTop(0).
			PaddingBottom(0)

	ListItemSelectedStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				Background(SelectedColor).
				PaddingLeft(1).
				PaddingRight(1).
				PaddingTop(0).
				PaddingBottom(0)

	// Item number style
	ItemNumberStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true).
			Width(4).
			Align(lipgloss.Right)

	// Command text style
	CommandTextStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			MarginLeft(1)

	// Empty state style
	EmptyStateStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Align(lipgloss.Center).
			Padding(2)
)

// RenderSearchComponent renders the two-column search bar
func RenderSearchComponent(prompt string, inputValue string, mode string, terminalWidth int) string {
	// Calculate widths
	usableWidth := terminalWidth - 4 // Account for margins

	// 85% for input, 15% for mode
	// Subtract 4 for borders (2 chars per panel)
	contentWidth := usableWidth - 2
	inputPanelWidth := int(float64(contentWidth) * 0.85)
	modePanelWidth := contentWidth - inputPanelWidth

	// Build left panel content
	promptText := PromptStyle.Render(prompt)
	leftContent := promptText + inputValue

	leftPanel := SearchInputBoxStyle.
		Width(inputPanelWidth).
		Render(leftContent)

	// Build right panel content
	modeText := fmt.Sprintf("[%s]", mode)
	rightPanel := ModeBadgeStyle.
		Width(modePanelWidth).
		Render(modeText)

	// Join panels horizontally
	searchBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	return SearchContainerStyle.Render(searchBar)
}

// RenderListComponent renders the command list with custom styling
// Shows maximum 10 items at a time, like fzf
func RenderListComponent(commands []history.Command, selectedIndex int, terminalWidth int, terminalHeight int) string {
	if len(commands) == 0 {
		emptyMessage := "No commands found"
		return ListContainerStyle.
			Width(terminalWidth - 4).
			Height(12).
			Render(EmptyStateStyle.Render(emptyMessage))
	}

	// Maximum 10 items visible at once
	const maxVisibleItems = 10
	
	var items []string
	var startIndex, endIndex int
	
	// Determine which items to show based on selected index
	if len(commands) <= maxVisibleItems {
		// Show all items if we have less than 10
		startIndex = 0
		endIndex = len(commands)
	} else {
		// Show sliding window of 10 items centered on selected item
		if selectedIndex < maxVisibleItems/2 {
			startIndex = 0
			endIndex = maxVisibleItems
		} else if selectedIndex > len(commands)-maxVisibleItems/2 {
			startIndex = len(commands) - maxVisibleItems
			endIndex = len(commands)
		} else {
			startIndex = selectedIndex - maxVisibleItems/2
			endIndex = selectedIndex + maxVisibleItems/2
		}
	}

	// Render visible items
	for i := startIndex; i < endIndex && i < len(commands); i++ {
		cmd := commands[i]
		isSelected := i == selectedIndex
		
		// Render item number
		itemNumber := ItemNumberStyle.Render(fmt.Sprintf("%d", cmd.Index))
		
		// Render command text
		commandText := CommandTextStyle.Render(cmd.Text)
		
		// Combine number and command
		itemContent := lipgloss.JoinHorizontal(lipgloss.Top, itemNumber, commandText)
		
		// Apply appropriate style based on selection
		if isSelected {
			items = append(items, ListItemSelectedStyle.Render(itemContent))
		} else {
			items = append(items, ListItemStyle.Render(itemContent))
		}
	}

	// Join all items vertically
	listContent := strings.Join(items, "\n")
	
	// Render the container with fixed height (10 items + borders)
	return ListContainerStyle.
		Width(terminalWidth - 4).
		Height(12).
		Render(listContent)
}
