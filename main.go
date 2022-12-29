package main

import (
	"fmt"
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/KaiAragaki/mimir-cli/tui"

	//"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"os"
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
