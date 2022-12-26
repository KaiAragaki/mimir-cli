package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func openDb() *gorm.DB {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/mimir??charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MSUSER"),
		os.Getenv("MSPASSWORD"),
		os.Getenv("MSHOST"),
		os.Getenv("MSPORT"),
	)

	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {

	db := openDb()

	fmt.Println(db)
	items := []list.Item{
		item{title: "one", desc: "list one"},
		item{title: "two", desc: "list two"},
		item{title: "three", desc: "list three"},
		item{title: "four", desc: "list four"},
	}
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Actions"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
