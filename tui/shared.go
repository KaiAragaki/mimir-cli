package tui

import (
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

// --- STRUCTS ---

// Entry
type Entry struct {
	template string  // Holds the print template
	fields   []field // The fields
	focused  int     // Which field is focused
	ok       bool    // Are all entries valid?
	repo     *gorm.DB
	subErr   string // What error (if any) came from submitting to DB?
}

// Field - a single unit of an entry
type field struct {
	input  textarea.Model
	hasErr bool
	errMsg string
	vfun   func(s string) (string, bool)
}

// Sensible defaults for fields
func NewDefaultField() field {
	ta := textarea.New()
	ta.FocusedStyle = textAreaFocusedStyle
	ta.ShowLineNumbers = false
	ta.Prompt = " "
	ta.BlurredStyle = textAreaBlurredStyle

	return field{
		input:  ta,
		hasErr: true,
		errMsg: "",
		vfun:   func(s string) (string, bool) { return "", false },
	}
}

// --- STYLING ---

const (
	white     = lipgloss.Color("#FFFFFF")
	purple    = lipgloss.Color("#7f12c7")
	darkGray  = lipgloss.Color("#767676")
	red       = lipgloss.Color("#FF0000")
	green     = lipgloss.Color("#00FF00")
	lightBlue = lipgloss.Color("#5C8DFF")
	blue      = lipgloss.Color("#3772FF")
	yellow    = lipgloss.Color("#FDCA40")
	black     = lipgloss.Color("#000000")
)

var (
	activeInputStyle     = lipgloss.NewStyle().Foreground(white).Background(purple)
	inactiveInputStyle   = lipgloss.NewStyle().Foreground(purple)
	continueStyle        = lipgloss.NewStyle().Foreground(darkGray)
	cursorStyle          = lipgloss.NewStyle().Foreground(white)
	cursorLineStyle      = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("230"))
	errorStyle           = lipgloss.NewStyle().Foreground(red)
	okStyle              = lipgloss.NewStyle().Foreground(green)
	textAreaFocusedStyle = textarea.Style{
		Base: lipgloss.
			NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(yellow).
			BorderLeft(true).
			Foreground(yellow),
	}
	textAreaBlurredStyle = textarea.Style{
		Base: lipgloss.
			NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(black).
			BorderLeft(true).
			Foreground(yellow),
	}
)

func newTextarea() textarea.Model {
	t := textarea.New()
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

func noFieldHasError(c Entry) bool {
	for _, v := range c.fields {
		if v.hasErr {
			return false
		}
	}
	return true
}

func getEntryStatus(c Entry) string {
	if noFieldHasError(c) {
		return okStyle.Render("Lookin' good!")
	} else {
		return errorStyle.Render("Entry not ready to be submitted.")
	}
}
