package tui

import (
	"fmt"

	"github.com/KaiAragaki/mimir-cli/db"
	"github.com/KaiAragaki/mimir-cli/shared"
	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/gorm"
)

// Field names ------
const (
	bcCellName = iota
	bcSeedConc
	bcSeedVol
	bcPlateFormFactor
	bcGrowthDuration
)

// Structures -------
type BaseCondition struct {
	Entry
}

func InitBaseCondition(findMode bool) tea.Model {
	inputs := make([]field, 5)
	for i := range inputs {
		inputs[i] = NewDefaultField()
		inputs[i].input.SetHeight(1)
	}

	inputs[bcCellName].displayName = "Cell Name"
	inputs[bcCellName].input.Focus()
	inputs[bcCellName].input.Placeholder = "j82"
	inputs[bcCellName].vfuns = append(
		inputs[bcCellName].vfuns,
		valIsBlank,
		valIsntLcAndNum,
	)

	inputs[bcSeedConc].displayName = "Seeding Concentration"
	inputs[bcSeedConc].input.Placeholder = "Number of cells/vol"
	inputs[bcSeedConc].vfuns = append(
		inputs[bcSeedConc].vfuns,
		valIsBlank,
		valMultiSlash,
		valStartsWithChar,
	)

	inputs[bcSeedVol].displayName = "Medium Volume"
	inputs[bcSeedVol].input.Placeholder = "Volume in well/flask/etc."
	inputs[bcSeedVol].vfuns = inputs[bcSeedConc].vfuns

	inputs[bcPlateFormFactor].displayName = "Plate Form Factor"
	inputs[bcPlateFormFactor].input.Placeholder = "6well"
	inputs[bcPlateFormFactor].vfuns = append(
		inputs[bcPlateFormFactor].vfuns,
		valIsBlank,
	)

	inputs[bcGrowthDuration].displayName = "Growth Time"
	inputs[bcGrowthDuration].input.Placeholder = "How long the cells grew until endpoint"
	inputs[bcGrowthDuration].vfuns = append(
		inputs[bcGrowthDuration].vfuns,
		valIsBlank,
		valIsntTimeUnit,
		valRepeatLetters,
		valStartsWithChar,
		valNoTimeUnit,
	)

	e := BaseCondition{
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

func (bc BaseCondition) Init() tea.Cmd {
	return nil
}

func (bc BaseCondition) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(bc.fields))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			table := InitTable(shared.Table)
			return table.Update(table)
		case tea.KeyTab:
			bc.focused = (bc.focused + 1) % len(bc.fields)
		case tea.KeyShiftTab:
			if bc.focused > 0 {
				bc.focused--
			}
		case tea.KeyEnter:
			// Don't newline in fields that are just 1 line tall
			// It's confusing!
			if bc.fields[bc.focused].input.Height() == 1 {
				return bc, nil
			}
		case tea.KeyCtrlS:
			if noFieldHasError(bc.Entry) {
				// TODO implement generalized makeEntry
				entry := makeBaseCondition(bc.Entry)
				err := bc.repo.Create(&entry).Error
				if err != nil {
					bc.subErr = errorStyle.Render(err.Error())
				} else {
					bc.subErr = okStyle.Render("Submitted!")
				}
			}
		default:
			// Only keep around submission errors
			// until the next key that could possibly fix it is pressed
			bc.subErr = ""
		}
		// Unfocus all inputs, then...
		for i := range bc.fields {
			bc.fields[i].input.Blur()
		}
		// Focus just the one
		bc.fields[bc.focused].input.Focus()
	}

	for i := range bc.fields {
		bc.fields[i].input, cmds[i] = bc.fields[i].input.Update(msg)
	}

	return bc, nil
}

func (bc BaseCondition) View() string {
	Validate(&bc.Entry)
	var out, header, err string
	for i, v := range bc.fields {
		if i == bc.focused {
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
		titleStyle.Render(" Add a Base Condition ") + "\n" +
			out +
			getEntryStatus(bc.Entry) + "\n\n" +
			bc.subErr + "\n\n")
}

func makeBaseCondition(e Entry) db.BaseCondition {
	var cell []db.Cell
	shared.DB.Where("cell_name = ?", e.fields[bcCellName].input.Value()).First(&cell)
	amt, _ := parseUnits(e.fields[bcSeedConc].input.Value())
	amtVol, _ := parseUnits(e.fields[bcSeedVol].input.Value())
	// Look up Cell by name here. If you don't find it, ask user if they want to make one.
	// If it does exist, fill out the info and store it in Cell.

	return db.BaseCondition{
		Model:                gorm.Model{},
		Cells:                cell,
		SeedingConcentration: amt,
		SeedingVolume:        amtVol,
		PlateFormFactor:      e.fields[bcPlateFormFactor].input.Value(),
		GrowthDuration:       parseTime(e.fields[bcGrowthDuration].input.Value()),
	}
}
