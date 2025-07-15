package repl

import "fmt"

func (m model) PreviewView() string {
	s := ""

	for i, v := range m.articles {
		if m.cursor == i {
			s += "> "
		}
		s += fmt.Sprintf("[ ] %d. %s \n", i+1, v.Title)
	}

	s += "\n \n"
	s += "l - enter article \n"
	s += "d - delete article \n"
	s += "q - quit \n"

	// s += fmt.Sprintf("C[%d]\n", m.cursor)

	return s
}
