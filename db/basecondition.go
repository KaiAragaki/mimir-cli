package db

import "gorm.io/gorm"

type BaseCondition struct {
	gorm.Model
	Cells                []Cell  `gorm:"many2many:basecondition_cells;"`
	SeedingConcentration float32 // cells/L
	SeedingVolume        float32 // SI
	PlateFormFactor      string
	GrowthDuration       int32 // time in seconds
}

func (r *Repo) AddBaseCondition(c *BaseCondition) {
	r.DB.Create(&c)
}
