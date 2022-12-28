package tui

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"

	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"
)

/*
 * Ideally I want something that looks like this:
 *
 * Cell Name:   _________  |
 * Parent Name: _________  |
 * Modifier:    _________  |
 *
 * That can become this when entries are formed:
 *
 * Cell Name:   UM-UC 15   | x Only numbers, underscores, and lowercase letters allowed
 * Parent Name: umuc15     | v
 * Modifier:    lorem ip...| v
 *
 * It should only submit when all checks are good, and should make sure there isn't a duplicate
 *
 * Editing should probably happen in a separate view
 * Errors may best be viewed in a separate view - just the indicator shows up,
 * but expand to see the complaints why.
 *
 * While some modifications could be made in post (convert to lowercase etc.),
 * I think it's a better idea to either make the user convert it, or convert it
 * in front of the user's eyes (when leaving the field). Will try the first,
 * and if it's too much of a pain, try the second
 *
 * May want some dummy text in the fields to give a good example
 *
 * Could asynchronously check to see if Cell Name exists in the DB while
 * the user adds more fields (use spinner to show checking in progress)
 */

type errMsg error

// Field names
const (
	cellName = iota
	parentName
	modifier
)

type Cell struct {
	template *template.Template // Holds the print template
	field    []field            // The fields
	focused  int                // Which field is focused
}

type field struct {
	input textinput.Model
	err   error
}

func InitCell() tea.Model {
	inputs := make([]field, 3)

	inputs[cellName].input = textinput.New()
	inputs[cellName].input.Focus()
	inputs[cellName].input.Validate = cellNameValidator

	inputs[parentName].input = textinput.New()
	inputs[parentName].input.Validate = cellNameValidator // A parent IS a cell so should be under identical strictures
	// It might also be useful to check if the parent is equal to the child. If so, it should be blank.

	inputs[modifier].input = textinput.New()
	inputs[modifier].input.Width = 20

	const cellTemplate = `
Add a cell entry:

  Cell Name {{ .CellName }}
{{ .CellNameError }}
Parent Name {{ .ParentName }}
{{ .ParentNameError }}
   Modifier {{ .Modifier }}
`

	template, err := template.New("err").Parse(cellTemplate)

	if err != nil {
		fmt.Printf("Error in templating: %v", err)
	}

	return Cell{
		template: template,
		field:    inputs,
		focused:  0,
	}
}

// ----- VALIDATORS -----
// BUG: Currently validators are blocking - so if something makes them upset,
// they prevent additional input.
func cellNameValidator(s string) error {
	//lcAndNum := regexp.MustCompile("^[a-z0-9]*$")
	//if !lcAndNum.MatchString(s) {
	//	return fmt.Errorf("May only include numbers and lowercase letters")
	//}

	return nil
}

// HACK Returns a string and a bool (for checking later)
func cellNameValidatorString(s string) (string, bool) {
	lcAndNum := regexp.MustCompile("^[a-z0-9]*$")
	if !lcAndNum.MatchString(s) {
		return "May only include numbers and lowercase letters", false
	}

	return "", true
}
func (c Cell) Init() tea.Cmd {
	return textinput.Blink // is this needed?
}

func (c Cell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(c.field))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			table := InitTable(shared.Table)
			return table.Update(table)
		case tea.KeyTab, tea.KeyEnter, tea.KeyDown:
			c.focused = (c.focused + 1) % len(c.field)
		case tea.KeyShiftTab, tea.KeyUp:
			if c.focused > 0 {
				c.focused--
			}
		}
		// Unfocus all inputs, then...
		for i := range c.field {
			c.field[i].input.Blur()
		}
		// Focus just the one
		c.field[c.focused].input.Focus()
	}

	for i := range c.field {
		c.field[i].input, cmds[i] = c.field[i].input.Update(msg)
	}

	return c, nil
}

func (c Cell) View() string {
	viewBuffer := &bytes.Buffer{}

	cellNameError, _ := cellNameValidatorString(c.field[cellName].input.Value())
	parentNameError, _ := cellNameValidatorString(c.field[parentName].input.Value())

	err := c.template.Execute(viewBuffer, map[string]interface{}{
		"CellName":        c.field[cellName].input.View(),
		"CellNameError":   errorStyle.Render(cellNameError),
		"ParentName":      c.field[parentName].input.View(),
		"ParentNameError": errorStyle.Render(parentNameError),
		"Modifier":        c.field[modifier].input.View(),
	})
	if err != nil {
		fmt.Println("Problem!!")
	}

	return viewBuffer.String()
}
