package repl

import "fmt"

func (m model) PreviewView() string {
	s := ""

	for i, v := range m.articles {
		s += fmt.Sprintf("[ ] %d. %s \n", i+1, v.Title)
	}

	var footer []string
	footer = append(footer, "l - enter article | d - delete article | q - quit")
	footer = append(footer, fmt.Sprintf("c:[%d] h:[%d]", m.cursor, m.sizes.height))

	visible := m.renderOnlyVisible(s, footer)

	return visible
}
