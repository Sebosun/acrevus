package repl

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) ViewArticle() string {
	s := ""

	// art := m.articles[m.selected]

	s += m.articleRawHtml

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
