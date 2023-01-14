package db

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"gorm.io/gorm"
)

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

func MakeBaseConditionTable() table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Cell Name", Width: 15},
		{Title: "Parent Name", Width: 15}, // In practice this can probably be removed - this is just for testing
		{Title: "Seeding Conc.", Width: 15},
		{Title: "Seeding Vol.", Width: 15},
		{Title: "Plate Form Factor", Width: 25},
		{Title: "Growth Time", Width: 15},
	}
	return table.New(table.WithColumns(columns))
}

func (bc BaseCondition) TableLineFromEntry() []table.Row {
	// Since each base condition can have multiple cells,
	// must loop through them?
	var rows []table.Row
	for _, v := range bc.Cells {
		row := []string{
			strconv.FormatUint(uint64(bc.ID), 10),
			v.CellName,
			v.ParentName,
			strconv.FormatFloat(float64(bc.SeedingConcentration), 'E', 4, 32),
			strconv.FormatFloat(float64(bc.SeedingVolume), 'E', 4, 32),
			bc.PlateFormFactor,
			strconv.FormatFloat(float64(bc.GrowthDuration), 'E', 4, 32),
		}
		rows = append(rows, row)
	}
	return rows
}
