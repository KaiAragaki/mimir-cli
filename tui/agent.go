package tui

import (
	"fmt"

	"github.com/KaiAragaki/mimir-cli/db"
	"github.com/KaiAragaki/mimir-cli/shared"
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

func InitAgent() tea.Model {
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
		valIsBlank,
		valIsntLcNumUnderDash,
	)

	inputs[amountWithUnits].displayName = "Amount with Units"
	inputs[amountWithUnits].input.Placeholder = "100nM"
	inputs[amountWithUnits].vfuns = append(
		inputs[amountWithUnits].vfuns,
		valIsBlank,
		valMultiSlash,
		valStartsWithChar,
	)

	inputs[agentDuration].displayName = "Agent Duration"
	inputs[amountWithUnits].input.Placeholder = "1d2m5s"
	inputs[agentDuration].vfuns = append(
		inputs[agentDuration].vfuns,
		valIsBlank,
		valIsntTimeUnit,
		valRepeatLetters,
		valStartsWithChar,
		valNoTimeUnit,
	)

	inputs[agentStartSincePlate].displayName = "Agent Start Since Plating"
	inputs[amountWithUnits].input.Placeholder = "1d"
	inputs[agentStartSincePlate].vfuns = inputs[agentDuration].vfuns

	e := Agent{
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

func (a Agent) Init() tea.Cmd {
	return nil
}

func (a Agent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(a.fields))
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
				// TODO implement generalized makeEntry
				entry := makeAgent(a.Entry)
				err := a.repo.Create(&entry).Error
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

	return a, nil
}

func (a Agent) View() string {
	Validate(&a.Entry)
	var out, header, err string
	for i, v := range a.fields {
		if i == a.focused {
			header = activeHeaderStyle.Render(" " + v.displayName)
		} else {
			header = headerStyle.Render(" " + v.displayName)
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

	return docStyle.Render(
		titleStyle.Render(" Add an Agent ") + "\n" +
			out +
			getEntryStatus(a.Entry) + "\n\n" +
			a.subErr + "\n\n")
}

func makeAgent(e Entry) db.Agent {
	amt, amtUnits := parseUnits(e.fields[amountWithUnits].input.Value())
	return db.Agent{
		Model:                gorm.Model{},
		AgentName:            e.fields[agentName].input.Value(),
		Amount:               amt,
		AmountUnits:          amtUnits,
		AgentDuration:        parseTime(e.fields[agentDuration].input.Value()),
		AgentStartSincePlate: parseTime(e.fields[agentStartSincePlate].input.Value()),
	}
}
