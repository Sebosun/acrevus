package repl

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	// The header
	if m.isReading {
		return m.displayArticle()
	}

	return m.displayOptions()
}

func (m model) displayArticle() string {
	s := ""

	art := m.articles[m.selected]

	s += art.RawText

	// The footer
	s += fmt.Sprintf("Width: %d, Height: %d", m.sizes.width, m.sizes.height)
	s += "\nPress q to quit.\n"
	s += fmt.Sprintf("%d - [%d]", m.cursor, m.sizes.height)

	style := lipgloss.NewStyle().
		Bold(false).
		Width(int(m.sizes.width/2)).
		Padding(1, 2).
		Foreground(lipgloss.Color("63"))

	stylized := style.Render(s)
	visible := m.centerText(m.renderOnlyVisible(stylized))

	return visible
}

func (m model) displayOptions() string {
	s := ""

	for i, v := range m.articles {
		if m.cursor == i {
			s += "> "
		}
		s += fmt.Sprintf("[ ] %d. %s \n", i+1, v.Title)
	}
	// s += fmt.Sprintf("%v %v", len(m.articles), m.cursor)
	return s
}

func (m model) centerText(rendered string) string {
	lines := strings.Split(rendered, "\n")
	var acc []string

	for _, v := range lines {
		innerAcc := ""
		for i := 0; i < m.sizes.width/(2*2); i++ {
			innerAcc += " "
		}
		innerAcc += v
		acc = append(acc, innerAcc)
	}

	return strings.Join(acc, "\n")
}

func (m model) renderOnlyVisible(rendered string) string {
	lines := strings.Split(rendered, "\n")
	var acc []string

	const padding int = 0

	for y, v := range lines {
		isIBigger := y >= m.cursor-padding
		isEnd := y < m.sizes.height+m.cursor-padding
		if isIBigger && isEnd {
			if y == m.cursor {
				acc = append(acc, ">"+v[1:])
			} else {
				acc = append(acc, v)
			}
		}
	}

	return strings.Join(acc, "\n")
}
