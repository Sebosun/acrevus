package repl

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	// The header
	if m.isReading {
		return m.ViewArticle()
	}

	return m.PreviewView()
}

func (m model) ViewArticle() string {
	s := ""

	art := m.articles[m.selected]

	s += art.RawText

	// The footer
	var footer []string

	footer = append(footer, fmt.Sprintf("\t \t \t  Width: %d, Height: %d \t     Press q to quit", m.sizes.width, m.sizes.height))
	footer = append(footer, fmt.Sprintf("c:[%d] h:[%d] S:[%d]", m.cursor, m.sizes.height, m.selected))

	style := lipgloss.NewStyle().
		Bold(false).
		Width(int(m.sizes.width/2)).
		Padding(1, 2)

	stylized := style.Render(s)
	visible := m.centerText(m.renderOnlyVisible(stylized, footer))

	return visible
}

func (m model) PreviewView() string {
	s := ""

	for i, v := range m.articles {
		if m.cursor == i {
			s += "> "
		}
		s += fmt.Sprintf("[ ] %d. %s \n", i+1, v.Title)
	}

	// s += fmt.Sprintf("C[%d]\n", m.cursor)

	return s
}
