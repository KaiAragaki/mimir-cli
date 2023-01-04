package tui

import (
	"fmt"

	"github.com/KaiAragaki/mimir-cli/db"
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/gorm"
)

// Field names -----
const (
	agentName = iota
	amountWithUnits
	agentDuration
	agentStartSincePlate
)

// Structures -----

type Agent struct {
	Entry
}

func InitAgent(findMode bool) tea.Model {
	inputs := make([]field, 4)
	for i := range inputs {
		inputs[i] = NewDefaultField()
		inputs[i].input.SetHeight(1)
	}

	inputs[agentName].displayName = "Agent Name"
	inputs[agentName].input.Focus()
	inputs[agentName].input.Placeholder = "saracatinib"
	inputs[agentName].vfuns = append(
		inputs[agentName].vfuns,
		valIsntLcNumUnderDash,
	)

	if !findMode {
		inputs[agentName].vfuns = append(inputs[agentName].vfuns, valIsBlank)
	}

	inputs[amountWithUnits].displayName = "Amount with Units"
	inputs[amountWithUnits].input.Placeholder = "100nM"
	inputs[amountWithUnits].vfuns = append(
		inputs[amountWithUnits].vfuns,
		valIsBlank,
		valMultiSlash,
		valStartsWithChar,
	)

	if !findMode {
		inputs[amountWithUnits].vfuns = append(inputs[amountWithUnits].vfuns, valIsBlank)
	}

	inputs[agentDuration].displayName = "Agent Duration"
	inputs[agentDuration].input.Placeholder = "1d2m5s"
	inputs[agentDuration].vfuns = append(
		inputs[agentDuration].vfuns,
		valIsBlank,
		valIsntTimeUnit,
		valRepeatLetters,
		valStartsWithChar,
		valNoTimeUnit,
	)

	if !findMode {
		inputs[agentDuration].vfuns = append(inputs[agentDuration].vfuns, valIsBlank)
	}

	inputs[agentStartSincePlate].displayName = "Agent Start Since Plating"
	inputs[agentStartSincePlate].input.Placeholder = "1d"
	inputs[agentStartSincePlate].vfuns = inputs[agentDuration].vfuns

	resTable := db.MakeAgentTable()

	e := Agent{
		Entry: Entry{
			fields:   inputs,
			focused:  0,
			ok:       false,
			repo:     shared.DB,
			subErr:   "",
			findMode: findMode,
			res:      resTable,
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

func (a Agent) Init() tea.Cmd {
	return nil
}

func (a Agent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(a.fields))

	entry := a.makeDbEntry()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			table := InitTable(shared.Table)
			return table.Update(table)
		case tea.KeyTab:
			a.focused = (a.focused + 1) % len(a.fields)
		case tea.KeyShiftTab:
			if a.focused > 0 {
				a.focused--
			}
		case tea.KeyEnter:
			// Don't newline in fields that are just 1 line tall
			// It's confusing!
			if a.fields[a.focused].input.Height() == 1 {
				return a, nil
			}
		case tea.KeyCtrlS:
			if noFieldHasError(a.Entry) {
				var err error
				if !a.findMode {
					err = a.repo.Create(&entry).Error
				}
				if err != nil {
					a.subErr = errorStyle.Render(err.Error())
				} else {
					a.subErr = okStyle.Render("Submitted!")
				}
			}
		default:
			// Only keep around submission errors
			// until the next key that could possibly fix it is pressed
			a.subErr = ""
		}
		// Unfocus all inputs, then...
		for i := range a.fields {
			a.fields[i].input.Blur()
		}
		// Focus just the one
		a.fields[a.focused].input.Focus()
	}

	for i := range a.fields {
		a.fields[i].input, cmds[i] = a.fields[i].input.Update(msg)
	}

	a = a.makeResTable(entry).(Agent)

	return a, nil
}

func (a Agent) View() string {
	Validate(&a.Entry)
	entry := a.makeDbEntry()
	a = a.makeResTable(entry).(Agent)
	var out, header, err, action string
	for i, v := range a.fields {
		if i == a.focused {
			header = activeHeaderStyle.Render(" " + v.displayName)
		} else {
			header = headerStyle.Render(" " + v.displayName)
		}

		if !a.findMode {
			if v.hasErr {
				err = errorStyle.Render(v.errMsg)
			} else {
				err = okStyle.Render("✓")
			}
		} else {
			err = ""
		}

		if v.hasErr {
			err = errorStyle.Render(v.errMsg)
		} else {
			if v.displayName == "Amount with Units" {
				parsedUnitsVal, parsedUnitsUnit := parseUnits(v.input.Value())
				out := fmt.Sprintf("%.5v %s", parsedUnitsVal, parsedUnitsUnit)
				err = okStyle.Render("✓ Will be converted to " + out)
			} else if v.displayName == "Agent Duration" || v.displayName == "Agent Start Since Plating" {
				parsedTime := parseTime(v.input.Value())
				out := fmt.Sprintf("%.5v", parsedTime)
				err = okStyle.Render("✓ Will be converted to " + out + " seconds")
			} else {
				err = okStyle.Render("✓")
			}

		}

		out = out + header + " " + err + "\n" +
			v.input.View() + "\n\n"
	}

	if a.findMode {
		action = "Find"
	} else {
		action = "Add"
	}

	return docStyle.Render(
		titleStyle.Render(" "+action+" an agent entry ") + "\n" +
			out +
			getEntryStatus(a.Entry) + "\n\n" +
			a.subErr + "\n\n" +
			a.res.View())
}

// UTILS ------------------
func (a Agent) makeDbEntry() db.Agent {
	amt, amtUnits := parseUnits(a.fields[amountWithUnits].input.Value())
	return db.Agent{
		Model:                gorm.Model{},
		AgentName:            a.fields[agentName].input.Value(),
		Amount:               amt,
		AmountUnits:          amtUnits,
		AgentDuration:        parseTime(a.fields[agentDuration].input.Value()),
		AgentStartSincePlate: parseTime(a.fields[agentStartSincePlate].input.Value()),
	}
}

func (a Agent) makeResTable(entry db.Agent) tea.Model {
	if a.findMode {
		var adb []db.Agent
		shared.DB.Where(entry).Limit(20).Find(&adb)
		var rows []table.Row
		for _, v := range adb {
			rows = append(rows, v.TableLineFromEntry())
		}
		a.res.SetRows(rows)
	} else {
		a.res.SetRows([]table.Row{entry.TableLineFromEntry()})
	}
	return a
}
