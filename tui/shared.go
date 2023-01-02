package tui

import (
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/gorm"
)

// --- STRUCTS ---

// Entry
type Entry struct {
	fields  []field // The fields
	focused int     // Which field is focused
	ok      bool    // Are all entries valid?
	repo    *gorm.DB
	subErr  string // What error (if any) came from submitting to DB?
}

type Enterable interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

// Field - a single unit of an entry
type field struct {
	displayName string // What the header of the field will be displayed as
	input       textarea.Model
	hasErr      bool
	errMsg      string
	vfuns       []func(s string) (string, bool)
}

// Sensible defaults for fields
func NewDefaultField() field {
	ta := textarea.New()
	ta.FocusedStyle = textAreaFocusedStyle
	ta.ShowLineNumbers = false
	ta.Prompt = " "
	ta.BlurredStyle = textAreaBlurredStyle
	//func(s string) (string, bool) { return "", false }
	var fns []func(s string) (string, bool)

	return field{
		input:  ta,
		hasErr: true,
		errMsg: "",
		vfuns:  fns,
	}
}

// FUNCTIONS ----------

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
	case "Agent":
		return InitAgent()
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

func Validate(c *Entry) {
	for i, v := range c.fields {
		for _, w := range v.vfuns {
			c.fields[i].errMsg, c.fields[i].hasErr = w(v.input.Value())
			if c.fields[i].hasErr {
				break
			}
		}
	}
}
