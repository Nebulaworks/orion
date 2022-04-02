package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#3FE0D0"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func checkbox(label string, checked bool, focused bool) string {
	var result = ""
	if checked {
		result = fmt.Sprintf("[x] %s", label)
	} else {
		result = fmt.Sprintf("[ ] %s", label)
	}

	if focused {
		return focusedStyle.Copy().Render(result)
	}
	return result
}
