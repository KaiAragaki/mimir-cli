package tui

import (
	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Table struct {
	list     list.Model
	selected string // which item is selected?
}

func InitTable(s string) tea.Model {
	items := []list.Item{
		item{title: "Cell", desc: "Those guys in flasks"},
		item{title: "Agent", desc: "Any experimental perturbation"},
		item{title: "Starting Condition", desc: "How many cells, for how long, in what flask"},
	}

	m := Table{
		list:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		selected: s,
	}
	m.list.Title = "Tables"
	return m
}

func (m Table) Init() tea.Cmd {
	return nil
}

func (m Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list.SetSize(shared.WindowSize.Width-1, shared.WindowSize.Height-2)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, shared.Keymap.Enter):
			shared.Table = m.list.SelectedItem().FilterValue()
			form := InitForm(shared.Table)
			return form.Update(shared.WindowSize)
		case key.Matches(msg, shared.Keymap.Back):
			action := InitAction()
			return action.Update(shared.WindowSize)
		}
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Table) View() string {
	return docStyle.Render(m.list.View())
}
