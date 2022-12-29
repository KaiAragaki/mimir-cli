package shared

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// VARIABLES ------
// These change, but they should be able to be accessed by everyone
var (
	DocStyle      = lipgloss.NewStyle()
	Action, Table string
	WindowSize    tea.WindowSizeMsg
	DB            = openDb()
)

// CONSTANTS ------
// These do NOT change

// FUNCTIONS ------
func openDb() *gorm.DB {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/mimir?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MSUSER"),
		os.Getenv("MSPASSWORD"),
		os.Getenv("MSHOST"),
		os.Getenv("MSPORT"),
	)

	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Keymaps
// Global keymaps that are centralized here
type keymap struct {
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
}

var Keymap = keymap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
