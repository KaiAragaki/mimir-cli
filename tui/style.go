package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

const (
	accent    = lipgloss.Color("#FDCA40")
	white     = lipgloss.Color("#FFFFFF")
	darkGray  = lipgloss.Color("#767676")
	vDarkGray = lipgloss.Color("#555555")
	green     = lipgloss.Color("#00FF00")
	black     = lipgloss.Color("#000000")
)

var (
	errorStyle       = lipgloss.NewStyle().Foreground(darkGray).Italic(true)
	okStyle          = lipgloss.NewStyle().Foreground(green)
	placeholderStyle = lipgloss.NewStyle().Foreground(vDarkGray)

	headerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(white).
			BorderLeft(true)

	activeHeaderStyle = headerStyle.Copy().
				Foreground(accent).
				Bold(true).
				BorderForeground(accent)

	textAreaFocusedStyle = textarea.Style{
		Base: activeHeaderStyle,
	}

	textAreaBlurredStyle = textarea.Style{
		Base: headerStyle,
	}

	titleStyle = lipgloss.NewStyle().
			Background(accent).
			Foreground(black).
			Margin(0, 2, 1, 2)

	docStyle = lipgloss.NewStyle().Margin(1)

	customTableStyle = table.Styles{
		Header:   headerStyle,
		Cell:     textAreaBlurredStyle.Text,
		Selected: textAreaBlurredStyle.Text,
	}
)

func newCustomListDelegate() list.ItemDelegate {
	customListDelegate := list.NewDefaultDelegate()
	customListDelegate.Styles.SelectedTitle.BorderForeground(accent)
	customListDelegate.Styles.SelectedTitle.Foreground(accent).Bold(true)
	customListDelegate.Styles.SelectedDesc.BorderLeftForeground(accent)
	customListDelegate.Styles.SelectedDesc.Foreground(accent)
	return customListDelegate
}
