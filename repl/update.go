package repl

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := m.alwaysUpdate(msg)
	if cmd != nil {
		return m, cmd
	}

	m, cmd = m.updateReading(msg)
	if cmd != nil {
		return m, cmd
	}

	m, cmd = m.updatePreview(msg)
	if cmd != nil {
		return m, cmd
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
		// These keys should exit the program.
		case "right", "l":
			m.cursor = 0
			m.selected = m.cursor
			m.isReading = true
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			m.cursor++

		}

	}
	return m, nil
}

func (m model) updatePreview(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) updateReading(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "left", "h":
			m.cursor = 0
			m.selected = m.cursor
			m.isReading = false

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.articles)-1 {
				m.cursor++
			}

		}
	}

	return m, nil
}
