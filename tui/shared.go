package tui

import (
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/gorm"

)

// --- STRUCTS ---

// Entry
type Entry struct {
	repo        *gorm.DB
	fields      []field     // The fields
	focused     int         // Which field is focused
	subErr      string      // What error (if any) came from submitting to DB?
	findMode    bool        // Should blank entries be ignored?
	res         table.Model // Stores find results or entry results
	entryStatus string      // What's good (or not) with the entry
	help  help.Model // Keystrokes for this view
}

// Field - a single unit of an entry
type field struct {
	input       textarea.Model
	displayName string                          // What the header of the field will be displayed as
	vfuns       []func(s string) (string, bool) // Validator functions
	hasErr      bool                            // Did the validator fail?
	errMsg      string                          // Error messages of valid
}

// Sensible defaults for fields
func NewDefaultField() field {
	ta := textarea.New()
	ta.FocusedStyle = textAreaFocusedStyle
	ta.ShowLineNumbers = false
	ta.Prompt = " "
	ta.BlurredStyle = textAreaBlurredStyle
	ta.BlurredStyle.Placeholder = placeholderStyle
	ta.FocusedStyle.Placeholder = placeholderStyle
	var fns []func(s string) (string, bool)

	return field{
		displayName: "",
		input:       ta,
		hasErr:      true,
		errMsg:      "",
		vfuns:       fns,
	}
}

// Sensible defaults for results table
// TODO: use DB entry style as structure, not mimir entry style
func NewDefaultTable(i []field) table.Model {
	columns := []table.Column{
		{Title: "id", Width: 4},
	}
	for _, v := range i {
		col := table.Column{Title: v.displayName, Width: len(v.displayName) + 4}
		columns = append(columns, col)
	}
	t := table.New(table.WithColumns(columns))
	return t
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
func InitForm(tableName string, findMode bool) tea.Model {
	switch tableName {
	case "Cell":
		return InitCell(findMode)
	case "Agent":
		return InitAgent(findMode)
	case "Base Condition":
		return InitBaseCondition(findMode)
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
	}
	return errorStyle.Render("Entry not ready to be submitted.")
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
