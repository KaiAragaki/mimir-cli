package db

import (
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

func (r *Repo) AddAgent(c *Agent) {
	r.DB.Create(&c)
}
