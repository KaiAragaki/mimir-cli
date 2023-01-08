package db

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
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

func MakeCellTable() table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Cell Name", Width: 15},
		{Title: "Parent Name", Width: 15},
		{Title: "Modifier", Width: 40},
	}
	t := table.New(table.WithColumns(columns))
	return t
}
