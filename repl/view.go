package repl

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	// The header
	s := "Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "1 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "2 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "3 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "4 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "5 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "6 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "7 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "8 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "9 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "10 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "11 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "12 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "13 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "14 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "15 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "16 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "17 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "18 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "19 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "20 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "21 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "22 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "23 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "24 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "25 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "26 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "27 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "28 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "29 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "30 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "31 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "32 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "33 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "34 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "35 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "36 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "37 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "38 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "39 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "40 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "41 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "42 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "43 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "44 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "45 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "46 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "47 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "48 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "49 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "50 Text text text hehe Text text text hehe Text text text hehe ?\n\n"
	s += "51 Text text text hehe Text text text hehe Text text text hehe ?\n\n"

	// The footer
	s += fmt.Sprintf("Width: %d, Height: %d", m.sizes.width, m.sizes.height)
	s += "\nPress q to quit.\n"
	s += fmt.Sprintf("%d - [%d]", m.cursor, m.sizes.height)

	style := lipgloss.NewStyle().
		Bold(true).
		Width(24).
		Padding(1, 2).
		Foreground(lipgloss.Color("63"))

	stylized := style.Render(s)
	visible := m.renderOnlyVisible(stylized)

	return visible
}

func (m model) renderOnlyVisible(rendered string) string {
	lines := strings.Split(rendered, "\n")
	var acc []string

	const padding int = 5

	for y, v := range lines {
		isIBigger := y >= m.cursor-padding
		isEnd := y < m.sizes.height+m.cursor-padding
		if isIBigger && isEnd {
			if y == m.cursor {
				acc = append(acc, ">"+v)
			} else {
				acc = append(acc, v)
			}
		}
	}

	return strings.Join(acc, "\n")
}
