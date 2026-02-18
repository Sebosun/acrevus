package repl

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := m.alwaysUpdate(msg)
	if cmd != nil {
		return m, cmd
	}

	switch m.state {
	case StateReading:
		m, cmd = m.updateReading(msg)
		if cmd != nil {
			return m, cmd
		}
	case StateMain:
		m, cmd = m.updatePreview(msg)
		if cmd != nil {
			return m, cmd
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) alwaysUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.sizes.width = msg.Width
		m.sizes.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}
