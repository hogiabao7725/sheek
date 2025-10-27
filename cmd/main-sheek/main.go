package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"sheek/internal/history"
	"sheek/internal/tui"
)

type teaModel tui.Model

func (m teaModel) Init() tea.Cmd { return nil }
func (m teaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	nm, cmd := tui.Update(msg, tui.Model(m))
	return teaModel(nm), cmd
}
func (m teaModel) View() string { return tui.View(tui.Model(m)) }

func main() {
	cmds, err := history.LoadAndParseZshHistory()
	if err != nil {
		log.Fatalf("failed to load history: %v", err)
	}

	model := tui.NewModel(cmds)

	// 💡 Dùng AltScreen để vẽ toàn màn hình, không đè prompt cũ
	p := tea.NewProgram(teaModel(model), tea.WithAltScreen())

	// 🧽 Chạy chương trình, sau khi Quit thì dọn sạch terminal
	if _, err := p.Run(); err != nil {
		log.Fatalf("error running program: %v", err)
	}

	// 💨 Clear màn hình sau khi chương trình kết thúc (fzf style)
	print("\033[H\033[2J")
}
