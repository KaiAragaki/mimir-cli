package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Submit     key.Binding
	FocusTable key.Binding
	Back       key.Binding
}

var FieldEntryKeyMap = KeyMap{
	Down: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "Down"),
	),
	Up: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("Shift + Tab", "Up"),
	),
	Submit: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("Ctrl + s", "Add to DB"),
	),
	FocusTable: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("Ctrl + t", "Toggle Table Focus"),
	),
	Back: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("Ctrl + c", "Back"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		FieldEntryKeyMap.Down,
		FieldEntryKeyMap.Up,
		FieldEntryKeyMap.Submit,
		FieldEntryKeyMap.FocusTable,
		FieldEntryKeyMap.Back,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return nil
}
