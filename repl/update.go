package repl

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.sizes.width = msg.Width
		m.sizes.height = msg.Height
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		//
		case "left", "h":
			m.cursor = 0
			m.selected = m.cursor
			m.isReading = false

		case "right", "l":
			m.cursor = 0
			m.selected = m.cursor
			m.isReading = true

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if !m.isReading {
				if m.cursor < len(m.articles)-1 {
					m.cursor++
				}
			} else {
				m.cursor++
			}

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) getArticleByIndex() {
}
