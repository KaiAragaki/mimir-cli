package cell

import (
	"gorm.io/gorm"
)

type Cell struct {
	gorm.Model
	CellName   string
	ParentName string
	Modifier   string
}

type Repo struct {
	DB *gorm.DB
}

func (r *Repo) AddCell(c *Cell) {
	r.DB.Create(&c)
}
