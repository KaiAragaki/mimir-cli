package shared

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// VARIABLES ------
// These change, but they should be able to be accessed by everyone
var (
	DocStyle      = lipgloss.NewStyle()
	Action, Table string
	WindowSize    tea.WindowSizeMsg
)

// CONSTANTS ------
// These do NOT change

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
