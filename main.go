package main

import (
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/KaiAragaki/mimir-cli/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	shared.Action = "Add"
	shared.Table = "Cell"
	m := tui.InitAction()

	p := tea.NewProgram(m)

	p.Run()
}
