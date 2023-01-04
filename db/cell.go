package db

import (
	"strconv"

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

func (c Cell) TableLineFromEntry() []string {
	return []string{
		strconv.FormatUint(uint64(c.ID), 10),
		c.CellName,
		c.ParentName,
		c.Modifier,
	}
}
