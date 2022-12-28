package tui

import (
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- STYLING ---

const (
	white    = lipgloss.Color("#FFFFFF")
	purple   = lipgloss.Color("#7f12c7")
	darkGray = lipgloss.Color("#767676")
	red      = lipgloss.Color("#FF0000")
)

var (
	activeInputStyle   = lipgloss.NewStyle().Foreground(white).Background(purple)
	inactiveInputStyle = lipgloss.NewStyle().Foreground(purple)
	continueStyle      = lipgloss.NewStyle().Foreground(darkGray)
	cursorStyle        = lipgloss.NewStyle().Foreground(white)
	cursorLineStyle    = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("230"))
	errorStyle         = lipgloss.NewStyle().Foreground(red)
)

func newTextInput() textinput.Model {
	t := textinput.New()
	t.CursorStyle = cursorStyle
	return t
}

// Both Action and Table share the same item structure, so it's defined here
type item struct {
	title, desc string
}

func (i item) Title() string {
	return i.title
}

func (i item) Description() string {
	return i.desc
}

func (i item) FilterValue() string {
	return i.title
}

// A function calls the correct Init* function based on the table name selected
// I'm sure there's a better way to do this (generics?) but I'm too dumb
func InitForm(tableName string) tea.Model {
	switch tableName {
	case "Cell":
		return InitCell()
	}
	return InitTable(shared.Table)
}
