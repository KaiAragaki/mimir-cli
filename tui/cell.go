package tui

import (
	"github.com/KaiAragaki/mimir-cli/db"
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

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

func InitCell(findMode bool) tea.Model {
	inputs := make([]field, 3)
	for i := range inputs {
		inputs[i] = NewDefaultField()
		inputs[i].input.SetHeight(1)
	}

	inputs[cellName].displayName = "Cell Name"
	inputs[cellName].input.Focus()
	inputs[cellName].input.Placeholder = "umuc6src54"
	inputs[cellName].vfuns = append(
		inputs[cellName].vfuns,
		valIsntLcNumUnder,
	)

	if !findMode {
		inputs[cellName].vfuns = append(inputs[cellName].vfuns, valIsBlank)
	}

	inputs[parentName].displayName = "Parent Name"
	inputs[parentName].input.Placeholder = "umuc6"
	inputs[parentName].vfuns = append(
		inputs[parentName].vfuns,
		valIsntLcNumUnder,
	)

	inputs[modifier].displayName = "Modifier"
	inputs[modifier].input.Placeholder = `Cells were transduced with...`
	inputs[modifier].input.SetHeight(5)
	inputs[modifier].hasErr = false

	resTable := db.MakeCellTable()

	e := Cell{
		Entry: Entry{
			repo:        shared.DB,
			fields:      inputs,
			focused:     0,
			subErr:      "",
			findMode:    findMode,
			res:         resTable,
			entryStatus: "",
			help:        help.New(),
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

	entry := c.makeDbEntry()
	// Format the entries for how they're ACTUALLY going to be searched for
	// and give a little preview in the table below

	c = c.makeResTable(entry).(Cell)

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
				var err error
				if !c.findMode {
					err = c.repo.Create(&entry).Error
				}
				if err != nil {
					c.subErr = errorStyle.Render(err.Error())
				} else {
					c.subErr = okStyle.Render("Submitted!")
				}
			}
		default:
			// Only keep around submission errors
			// until the next key that could possibly fix it is pressed
			c.entryStatus = getEntryStatus(c.Entry)
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

	c.makeResTable(entry)
	c.res.SetStyles(customTableStyle)
	return c, nil
}

func (c Cell) View() string {
	Validate(&c.Entry)
	entry := c.makeDbEntry()
	c = c.makeResTable(entry).(Cell)
	var out, header, err, action string
	for i, v := range c.fields {
		if i == c.focused {
			header = activeHeaderStyle.Render(" " + v.displayName)
		} else {
			header = headerStyle.Render(" " + v.displayName)
		}

		if !c.findMode {
			if v.hasErr {
				err = errorStyle.Render(v.errMsg)
			} else {
				err = okStyle.Render("✓")
			}
		} else {
			err = ""
		}

		out = out + header + " " + err + "\n" +
			v.input.View() + "\n\n"
	}

	if c.findMode {
		action = "Find"
	} else {
		action = "Add"
	}

	leftCol := out + getEntryStatus(c.Entry)

	return docStyle.Render(
		titleStyle.Render(" "+action+" a cell entry ") + "\n\n" +
			lipgloss.JoinHorizontal(0, wholeTableStyle.Render(leftCol), wholeTableStyle.Render(c.res.View())) + "\n\n" +
			c.help.View(FieldEntryKeyMap),
	)
}

// UTILS ------------------

// Constructor of db entry
func (c Cell) makeDbEntry() db.Cell {
	return db.Cell{
		Model:      gorm.Model{},
		CellName:   c.fields[cellName].input.Value(),
		ParentName: c.fields[parentName].input.Value(),
		Modifier:   c.fields[modifier].input.Value(),
	}
}

func (c Cell) makeResTable(entry db.Cell) tea.Model {
	if c.findMode {
		var cdb []db.Cell
		shared.DB.Where(entry).Limit(20).Find(&cdb)
		var rows []table.Row
		for _, v := range cdb {
			rows = append(rows, v.TableLineFromEntry())
		}
		c.res.SetRows(rows)
	} else {
		c.res.SetRows([]table.Row{entry.TableLineFromEntry()})
	}
	return c
}
