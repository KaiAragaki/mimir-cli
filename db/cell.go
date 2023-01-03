package db

import (
	"gorm.io/gorm"
)

type Cell struct {
	gorm.Model
	CellName   string `gorm:"size:255"`
	ParentName string
	Modifier   string
}

func (r *Repo) AddCell(c *Cell) {
	r.DB.Create(&c)
}
