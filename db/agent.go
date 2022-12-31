package db

import (
	"gorm.io/gorm"
)

type Agent struct {
	gorm.Model
	AgentName            string
	Concentration        float32
	ConcentrationUnits   string // SI
	AgentDuration        int32  // time in seconds (max: 68y)
	AgentStartSincePlate int32  // time in seconds
}

// ??interface for adding? so it can be general?
func (r *Repo) AddAgent(c *Agent) {
	r.DB.Create(&c)
}
