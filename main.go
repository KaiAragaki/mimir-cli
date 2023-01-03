package main

import (
	"fmt"
	"os"

	"github.com/KaiAragaki/mimir-cli/db"
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/KaiAragaki/mimir-cli/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	shared.Action = "Add"
	shared.Table = "Cell"
	m := tui.InitAction()
	migErr := shared.DB.AutoMigrate(&db.Agent{}, &db.Cell{})
	if migErr != nil {
		fmt.Println("There was an error migrating the database:", migErr)
		os.Exit(1)
	}
	migErr2 := shared.DB.AutoMigrate(&db.BaseCondition{})
	if migErr2 != nil {
		fmt.Println("There was an error migrating the database:", migErr2)
		os.Exit(1)
	}
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
