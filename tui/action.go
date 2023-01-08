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

func InitAction() tea.Model {
	items := []list.Item{
		item{title: "Add", desc: "Add an item"},
		item{title: "Find", desc: "Look up, edit, or delete an item"},
	}

	l := list.NewModel(items, newCustomListDelegate(), 0, 0)
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle.Padding(0, 1).Margin(0)

	m := Action{
		list:     l,
		selected: shared.Action,
	}
	m.list.Title = "Actions"
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
	case tea.WindowSizeMsg:
		shared.WindowSize = msg
		m.list.SetSize(shared.WindowSize.Width-1, shared.WindowSize.Height-2)
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Action) View() string {
	return docStyle.Render(m.list.View())
}
