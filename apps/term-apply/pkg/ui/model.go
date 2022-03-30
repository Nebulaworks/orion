package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nebulaworks/orion/apps/term-apply/pkg/applicant"
)

type Model struct {
	focusIndex int
	choice     int
	checkBoxes []string
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
	Submitted  bool
	sub        chan responseMsg // where we'll receive activity notifications
	response   string
	appMgr     *applicant.ApplicantManager
	userID     string
}

func InitialModel(am *applicant.ApplicantManager, user string) Model {
	m := Model{
		inputs:     make([]textinput.Model, 2),
		checkBoxes: make([]string, 2),
		sub:        make(chan responseMsg),
		appMgr:     am,
		userID:     user,
		response:   "not found",
	}

	m.checkBoxes[0] = "Senior Software Engineer"
	m.checkBoxes[1] = "Software Engineer"

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Full Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Email"
			t.CharLimit = 64
		}

		m.inputs[i] = t
	}

	return m
}
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.listenForActivity(m.sub), // generate activity
		waitForActivity(m.sub),     // wait for activity
	)

}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case responseMsg:
		m.response = "received, thank you"
		return m, waitForActivity(m.sub) // wait for next event
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == (len(m.inputs)+len(m.checkBoxes)) {
				if err := m.appMgr.AddApplicant(
					m.userID,
					m.inputs[0].Value(),
					m.inputs[1].Value(),
					m.choice,
				); err == nil {
					m.Submitted = true
				} else {
					m.Submitted = false
				}
			} else if s == "enter" && m.focusIndex >= 2 {
				m.choice = m.focusIndex - 2
				m.focusIndex = m.focusIndex - 1
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > (len(m.inputs) + len(m.checkBoxes)) {
				// if we exceed the focus index just keep us where we were
				m.focusIndex = m.focusIndex - 1
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) + len(m.checkBoxes)
			}

			cmds := make([]tea.Cmd, (len(m.inputs) + len(m.checkBoxes)))
			for i := 0; i <= (len(m.inputs)+len(m.checkBoxes))-1; i++ {
				if i == m.focusIndex && i < 2 {
					// if re-editing info
					m.Submitted = false
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				} else if i < 2 {

					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}

			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	b.WriteRune('\n')
	b.WriteRune('\n')
	var focus = false
	for i := range m.checkBoxes {
		if m.focusIndex-2 == i {
			focus = true
		} else {
			focus = false
		}

		if m.choice == i {
			b.WriteString(checkbox(m.checkBoxes[i], true, focus))
		} else {
			b.WriteString(checkbox(m.checkBoxes[i], false, focus))
		}

		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs)+len(m.checkBoxes) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	if m.Submitted {
		b.WriteString(fmt.Sprintf("\n%s thank you for applying to Nebulaworks! We will follow up with you shortly via your email:%s regarding next steps\n", m.inputs[0].Value(), m.inputs[1].Value()))
	}
	b.WriteString(fmt.Sprintf("\n Resume status: %s \n\n", m.response))
	b.WriteString(helpStyle.Render("ctrl+c to exit"))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}
