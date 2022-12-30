package db

import (
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	CellName   string
	ParentName string
	Modifier   string
}

type Repo struct {
	DB *gorm.DB
}
