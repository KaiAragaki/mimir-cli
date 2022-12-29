package tui

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"

	"github.com/KaiAragaki/mimir-cli/cell"
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"
)

type errMsg error

// Field names ------
const (
	cellName = iota
	parentName
	modifier
)

// Define Structures ------
type Cell struct {
	template *template.Template // Holds the print template
	fields   []field            // The fields
	focused  int                // Which field is focused
	ok       bool               // Are all entries valid?
	repo     cell.Repo
}

type field struct {
	input  textinput.Model
	hasErr bool
	errMsg string
	vfun   func(s string) (string, bool)
}

func NewField() field {
	return field{
		input:  textinput.New(),
		hasErr: true,
		errMsg: "",
		vfun:   func(s string) (string, bool) { return "", false },
	}
}
func InitCell() tea.Model {
	inputs := make([]field, 3)
	for i := range inputs {
		inputs[i] = NewField()
	}

	inputs[cellName].input.Focus()
	inputs[cellName].vfun = cellNameValidatorString

	inputs[parentName].vfun = parentNameValidatorString

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
		fields:   inputs,
		focused:  0,
		ok:       false,
		repo:     cell.Repo{DB: shared.DB},
	}
}

// VALIDATORS -------------
// BUG: Currently validators are blocking - so if something makes them upset,
// they prevent additional input.
// HACK Returns a string and a bool (for checking later)
func cellNameValidatorString(s string) (string, bool) {
	lcAndNum := regexp.MustCompile("^[a-z0-9]*$")
	if !lcAndNum.MatchString(s) {
		return "May only include numbers and lowercase letters", true
	}

	if s == "" {
		return "Field must not be blank", true
	}

	return "", false
}

func parentNameValidatorString(s string) (string, bool) {
	lcAndNum := regexp.MustCompile("^[a-z0-9]*$")
	if !lcAndNum.MatchString(s) {
		return "May only include numbers and lowercase letters", true
	}

	return "", false
}

func updateErrors(c *Cell) {
	for i, v := range c.fields {
		c.fields[i].errMsg, c.fields[i].hasErr = v.vfun(v.input.Value())
	}
}

func (c Cell) Init() tea.Cmd {
	return textinput.Blink // is this needed?
}

func (c Cell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(c.fields))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			table := InitTable(shared.Table)
			return table.Update(table)
		case tea.KeyTab, tea.KeyDown:
			c.focused = (c.focused + 1) % len(c.fields)
		case tea.KeyShiftTab, tea.KeyUp:
			if c.focused > 0 {
				c.focused--
			}
		case tea.KeyEnter:
			// check if complete and ok here
			if noFieldHasError(c) {
				entry := makeCell(c)
				c.repo.AddCell(&entry)
			}
		}
		// Unfocus all inputs, then...
		for i := range c.fields {
			c.fields[i].input.Blur()
		}
		// Focus just the one
		c.fields[c.focused].input.Focus()
	}

	for i := range c.fields {
		c.fields[i].input, cmds[i] = c.fields[i].input.Update(msg)
	}

	return c, nil
}

func (c Cell) View() string {
	viewBuffer := &bytes.Buffer{} // Might cause misfiring with submission. We'll see.

	updateErrors(&c)
	err := c.template.Execute(viewBuffer, map[string]interface{}{
		"CellName":        c.fields[cellName].input.View(),
		"CellNameError":   errorStyle.Render(c.fields[cellName].errMsg),
		"ParentName":      c.fields[parentName].input.View(),
		"ParentNameError": errorStyle.Render(c.fields[parentName].errMsg),
		"Modifier":        c.fields[modifier].input.View(),
	})
	if err != nil {
		fmt.Println("Error creating template")
	}
	return viewBuffer.String()
}

// UTILS -----

func noFieldHasError(c Cell) bool {
	for _, v := range c.fields {
		if v.hasErr {
			return false
		}
	}
	return true
}

func makeCell(c Cell) cell.Cell {
	return cell.Cell{
		CellName:   c.fields[cellName].input.Value(),
		ParentName: c.fields[parentName].input.Value(),
		Modifier:   c.fields[modifier].input.Value(),
	}
}
