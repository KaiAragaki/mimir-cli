package tui

import (
	"regexp"

	"github.com/KaiAragaki/mimir-cli/cell"
	"github.com/KaiAragaki/mimir-cli/shared"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

// Field names ------
const (
	cellName = iota
	parentName
	modifier
)

// Define Structures ------
func InitCell() tea.Model {
	inputs := make([]field, 3)
	for i := range inputs {
		inputs[i] = NewDefaultField()
		inputs[i].input.SetHeight(1)
	}

	inputs[cellName].displayName = "Cell Name"
	inputs[cellName].input.Focus()
	inputs[cellName].vfun = cellNameValidatorString

	inputs[parentName].displayName = "Parent Name"
	inputs[parentName].vfun = parentNameValidatorString

	inputs[modifier].displayName = "Modifier"
	inputs[modifier].input.SetHeight(5)

	const tmpl = `
Add a cell entry:

Cell Name
%s
%s
Parent Name
%s
%s
Modifier
%s

%s

%s
`
	return Entry{
		template: tmpl,
		fields:   inputs,
		focused:  0,
		ok:       false,
		repo:     shared.DB,
		subErr:   "",
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

func updateErrors(c *Entry) {
	for i, v := range c.fields {
		c.fields[i].errMsg, c.fields[i].hasErr = v.vfun(v.input.Value())
	}
}

func (c Entry) Init() tea.Cmd {
	return nil
}

func (c Entry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(c.fields))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			table := InitTable(shared.Table)
			return table.Update(table)
		case tea.KeyTab:
			c.focused = (c.focused + 1) % len(c.fields)
		case tea.KeyShiftTab:
			if c.focused > 0 {
				c.focused--
			}
		case tea.KeyEnter:
			if noFieldHasError(c) {
				entry := makeCell(c)
				err := c.repo.Create(&entry).Error
				if err != nil {
					c.subErr = errorStyle.Render(err.Error())
				} else {
					c.subErr = okStyle.Render("Submitted!")
				}
			}
		default:
			c.subErr = ""
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

func (c Entry) View() string {

	updateErrors(&c) // Might cause misfiring with submission. We'll see.

	var out string

	for _, v := range c.fields {
		out = out + v.displayName + "\n" +
			v.input.View() + "\n" +
			errorStyle.Render(v.errMsg) + "\n"
	}

	return "Add a cell entry\n\n" +
		out + "\n\n" +
		getEntryStatus(c) + "\n\n" +
		c.subErr
}

// UTILS -----
func makeCell(c Entry) cell.Cell {
	return cell.Cell{
		CellName:   c.fields[cellName].input.Value(),
		ParentName: c.fields[parentName].input.Value(),
		Modifier:   c.fields[modifier].input.Value(),
	}
}
