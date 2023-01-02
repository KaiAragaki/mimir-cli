package tui

import (
	"github.com/KaiAragaki/mimir-cli/db"
	"github.com/KaiAragaki/mimir-cli/shared"
	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/gorm"
)

var debug string

type errMsg error

// Field names ------
const (
	cellName = iota
	parentName
	modifier
)

// Define Structures ------

type Cell struct {
	Entry
}

func InitCell() tea.Model {
	inputs := make([]field, 3)
	for i := range inputs {
		inputs[i] = NewDefaultField()
		inputs[i].input.SetHeight(1)
	}

	inputs[cellName].displayName = "Cell Name"
	inputs[cellName].input.Focus()
	inputs[cellName].vfuns = append(
		inputs[cellName].vfuns,
		valIsBlank,
		valIsntLcAndNum,
	)

	inputs[parentName].displayName = "Parent Name"
	inputs[parentName].vfuns = append(
		inputs[parentName].vfuns,
		valIsntLcAndNum,
	)

	inputs[modifier].displayName = "Modifier"
	inputs[modifier].input.SetHeight(5)
	inputs[modifier].hasErr = false

	e := Cell{
		Entry: Entry{
			fields:  inputs,
			focused: 0,
			ok:      false,
			repo:    shared.DB,
			subErr:  "",
		},
	}

	// Initialize all foci so there's no pop in
	for i := range e.fields {
		e.fields[i].input.Blur()
	}
	// Focus just the one
	e.fields[e.focused].input.Focus()

	return e
}

func (c Cell) Init() tea.Cmd {
	return nil
}

func (c Cell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// Don't newline in fields that are just 1 line tall
			// It's confusing!
			if c.fields[c.focused].input.Height() == 1 {
				return c, nil
			}
		case tea.KeyCtrlS:
			if noFieldHasError(c.Entry) {
				// TODO implement generalized makeEntry
				entry := makeCell(c.Entry)
				err := c.repo.Create(&entry).Error
				if err != nil {
					c.subErr = errorStyle.Render(err.Error())
				} else {
					c.subErr = okStyle.Render("Submitted!")
				}
			}
		default:
			// Only keep around submission errors
			// until the next key that could possibly fix it is pressed
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

func (c Cell) View() string {
	Validate(&c.Entry)
	var out, header, err string
	for i, v := range c.fields {
		if i == c.focused {
			header = activeHeaderStyle.Render(" " + v.displayName)
		} else {
			header = headerStyle.Render(" " + v.displayName)
		}

		if v.hasErr {
			err = errorStyle.Render(v.errMsg)
		} else {
			err = okStyle.Render("✓")
		}

		out = out + header + " " + err + "\n" +
			v.input.View() + "\n\n"
	}

	return docStyle.Render(
		titleStyle.Render(" Add a cell entry ") + "\n" +
			out +
			getEntryStatus(c.Entry) + "\n\n" +
			c.subErr + "\n\n")
}

// UTILS ------------------

// Constructor of db entry
func makeCell(c Entry) db.Cell {
	return db.Cell{
		Model:      gorm.Model{},
		CellName:   c.fields[cellName].input.Value(),
		ParentName: c.fields[parentName].input.Value(),
		Modifier:   c.fields[modifier].input.Value(),
	}
}
