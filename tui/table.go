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
		item{title: "Base Condition", desc: "How many cells, for how long, in what flask"},
	}

	l := list.NewModel(items, newCustomListDelegate(), 0, 0)
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle.Padding(0, 1).Margin(0)

	m := Table{
		list:     l,
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
			var model tea.Model
			if shared.Action == "Add" {
				model = InitForm(shared.Table, false)
			} else if shared.Action == "Find" {
				model = InitForm(shared.Table, true)
			}
			return model.Update(shared.WindowSize)
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
