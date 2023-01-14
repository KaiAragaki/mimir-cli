package tui

import (
	"fmt"

	"github.com/KaiAragaki/mimir-cli/db"
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		valIsntLcNumUnder,
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

	resTable := db.MakeBaseConditionTable()

	e := BaseCondition{
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

func (bc BaseCondition) Init() tea.Cmd {
	return nil
}

func (bc BaseCondition) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(bc.fields))

	entry := bc.makeDbEntry()

	bc = bc.makeResTable(entry).(BaseCondition)

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
				var err error
				if !bc.findMode {
					err = bc.repo.Create(&entry).Error
				}
				if err != nil {
					bc.subErr = errorStyle.Render(err.Error())
				} else {
					bc.subErr = okStyle.Render("Submitted!")
				}
			}
		default:
			// Only keep around submission errors
			// until the next key that could possibly fix it is pressed
			bc.entryStatus = getEntryStatus(bc.Entry)
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

	bc = bc.makeResTable(entry).(BaseCondition)
	bc.res.SetStyles(customTableStyle)
	return bc, nil
}

func (bc BaseCondition) View() string {
	Validate(&bc.Entry)
	entry := bc.makeDbEntry()
	bc = bc.makeResTable(entry).(BaseCondition)
	var out, header, err, action string
	for i, v := range bc.fields {
		if i == bc.focused {
			header = activeHeaderStyle.Render(" " + v.displayName)
		} else {
			header = " " + v.displayName
		}

		if !bc.findMode {
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

	if bc.findMode {
		action = "Find"
	} else {
		action = "Add"
	}

	leftCol := out + getEntryStatus(bc.Entry)

	return docStyle.Render(
		titleStyle.Render(" "+action+" a base condition entry ") + "\n\n" +
			lipgloss.JoinHorizontal(0, wholeTableStyle.Render(leftCol), wholeTableStyle.Render(bc.res.View())) + "\n\n" +
			bc.help.View(FieldEntryKeyMap),
	)
}

func (bc BaseCondition) makeDbEntry() db.BaseCondition {
	var cell []db.Cell
	shared.DB.Where("cell_name = ?", bc.fields[bcCellName].input.Value()).First(&cell)
	amt, _ := parseUnits(bc.fields[bcSeedConc].input.Value())
	amtVol, _ := parseUnits(bc.fields[bcSeedVol].input.Value())
	// Look up Cell by name here. If you don't find it, ask user if they want to make one.
	// If it does exist, fill out the info and store it in Cell.

	return db.BaseCondition{
		Model:                gorm.Model{},
		Cells:                cell,
		SeedingConcentration: amt,
		SeedingVolume:        amtVol,
		PlateFormFactor:      bc.fields[bcPlateFormFactor].input.Value(),
		GrowthDuration:       parseTime(bc.fields[bcGrowthDuration].input.Value()),
	}
}

func (bc BaseCondition) makeResTable(entry db.BaseCondition) tea.Model {
	if bc.findMode {
		var bcs []db.BaseCondition
		shared.DB.Where(entry).Preload("Cells").Find(&bcs)
		for _, v := range bcs {
			bc.res.SetRows(append(bc.res.Rows(), v.TableLineFromEntry()...))
		}
	} else {
		bc.res.SetRows(entry.TableLineFromEntry())
	}
	return bc
}
