package main

import (
	"fmt"
	"os"

	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/KaiAragaki/mimir-cli/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	shared.Action = "Add"
	shared.Table = "Cell"
	m := tui.InitAction()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
