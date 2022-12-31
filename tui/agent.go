package tui

import (
	"github.com/KaiAragaki/mimir-cli/shared"
	tea "github.com/charmbracelet/bubbletea"
)

// Field names -----
const (
	agentName = iota
	concentration
	concentrationUnits
	agentDuration
	agentStartSincePlate
)

// Structures -----
func InitAgent() tea.Model {
	inputs := make([]field, 5)
	for i := range inputs {
		inputs[i] = NewDefaultField()
		inputs[i].input.SetHeight(1)
	}

	inputs[agentName].displayName = "Agent Name"
	inputs[agentName].input.Focus()
	inputs[agentName].vfuns = append(
		inputs[agentName].vfuns,
		valIsBlank,
		valIsntLcNumUnderDash,
	)

	inputs[concentration].displayName = "Concentration"

	inputs[concentrationUnits].displayName = "Concentration Units"

	inputs[agentDuration].displayName = "Agent Duration"

	inputs[agentStartSincePlate].displayName = "Agent Start Since Plating"

	return Entry{
		fields:  inputs,
		focused: 0,
		ok:      false,
		repo:    shared.DB,
		subErr:  "",
	}
}

// UTILS -----
