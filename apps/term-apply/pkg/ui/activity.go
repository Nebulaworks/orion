package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type responseMsg struct{}

func (m *Model) listenForActivity(sub chan responseMsg) tea.Cmd {
	return func() tea.Msg {
		for {
			// only send the message the message when we first
			// find the resume
			if m.appMgr.HasResume(m.userID) {
				sub <- responseMsg{}
				break
			}
			time.Sleep(time.Second * time.Duration(5))
		}
		return nil
	}
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan responseMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}
