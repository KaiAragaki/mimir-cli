package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

const (
	white     = lipgloss.Color("#FFFFFF")
	purple    = lipgloss.Color("#7f12c7")
	darkGray  = lipgloss.Color("#767676")
	vDarkGray = lipgloss.Color("#555555")
	red       = lipgloss.Color("#FF0000")
	green     = lipgloss.Color("#00FF00")
	lightBlue = lipgloss.Color("#5C8DFF")
	blue      = lipgloss.Color("#3772FF")
	yellow    = lipgloss.Color("#FDCA40")
	black     = lipgloss.Color("#000000")
)

var (
	activeInputStyle   = lipgloss.NewStyle().Foreground(white).Background(purple)
	inactiveInputStyle = lipgloss.NewStyle().Foreground(purple)
	continueStyle      = lipgloss.NewStyle().Foreground(darkGray)
	cursorStyle        = lipgloss.NewStyle().Foreground(white)
	cursorLineStyle    = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("230"))
	errorStyle         = lipgloss.NewStyle().Foreground(darkGray).Italic(true)
	okStyle            = lipgloss.NewStyle().Foreground(green)
	placeholderStyle   = lipgloss.NewStyle().Foreground(vDarkGray)

	textAreaFocusedStyle = textarea.Style{
		Base: lipgloss.
			NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(yellow).
			BorderLeft(true).
			Foreground(yellow),
	}
	textAreaBlurredStyle = textarea.Style{
		Base: lipgloss.
			NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(white).
			BorderLeft(true).
			Foreground(white),
	}

	headerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(white).
			BorderLeft(true)

	activeHeaderStyle = headerStyle.Copy().
				Foreground(yellow).Bold(true).
				BorderForeground(yellow)

	titleStyle = lipgloss.NewStyle().
			Background(yellow).
			Foreground(black).
			Margin(0, 2, 3, 2)
	docStyle = lipgloss.NewStyle().Margin(1)
)
