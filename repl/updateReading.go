package repl

import tea "github.com/charmbracelet/bubbletea"

func (m model) updateReading(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "left", "h":
			m.cursor = 0
			m.selected = m.cursor
			m.state = StateMain

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			m.cursor++

		case "U":
			if m.cursor-m.sizes.height > 0 {
				m.cursor -= m.sizes.height
			} else {
				m.cursor = 0
			}

		case "D":
			m.cursor += m.sizes.height
		}
	}

	return m, nil
}
