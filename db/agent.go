package db

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model
	AgentName            string  `gorm:"size:255"`
	Amount               float32 // SI
	AmountUnits          string  `gorm:"size:255"` // SI
	AgentDuration        int32   // time in seconds
	AgentStartSincePlate int32   // time in seconds
}

func (r *Repo) AddAgent(a *Agent) {
	r.DB.Create(&a)
}

func (a Agent) TableLineFromEntry() []string {
	return []string{
		strconv.FormatUint(uint64(a.ID), 10),
		a.AgentName,
		strconv.FormatFloat(float64(a.Amount), 'E', 4, 32),
		a.AmountUnits,
		strconv.FormatInt(int64(a.AgentDuration), 10),
		strconv.FormatInt(int64(a.AgentStartSincePlate), 10),
	}
}

func MakeAgentTable() table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Agent Name", Width: 15},
		{Title: "Amount", Width: 10},
		{Title: "Amount Units", Width: 20},
		{Title: "Agent Duration", Width: 20},
		{Title: "Agent Start Since Plate", Width: 30},
	}
	t := table.New(table.WithColumns(columns), table.WithHeight(10))
	return t
}
