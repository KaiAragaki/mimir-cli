package tui

import (
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Action struct {
	list     list.Model
	selected string // which item is selected?
}

type item struct {
	title, desc string
}

func (i item) Title() string {
	return i.title
}

func (i item) Description() string {
	return i.desc
}

func (i item) FilterValue() string {
	return i.title
}

func InitAction() tea.Model {
	items := []list.Item{
		item{title: "Add", desc: "Add an item"},
		item{title: "Find", desc: "Look up, edit, or delete an item"},
	}
	m := Action{
		list:     list.NewModel(items, list.NewDefaultDelegate(), 8, 8),
		selected: shared.Action,
	}
	return m
}

func (m Action) Init() tea.Cmd {
	return nil
}

func (m Action) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, shared.Keymap.Enter):
			shared.Action = m.list.SelectedItem().FilterValue()
			table := InitTable(shared.Action)
			return table.Update(shared.WindowSize)
		}
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Action) View() string {
	return shared.DocStyle.Render(m.list.View())
}
