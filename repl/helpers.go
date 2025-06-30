package repl

import "strings"

func (m model) renderOnlyVisible(rendered string, footer []string) string {
	lines := strings.Split(rendered, "\n")
	footerLines := parseFooter(footer)

	var acc []string
	footLen := len(footerLines)

	for i, v := range lines {
		isBeforeView := i >= m.cursor
		isAfterView := i < m.sizes.height+m.cursor-footLen

		if isBeforeView && isAfterView {
			if i == m.cursor {
				acc = append(acc, ">"+v[1:])
			} else {
				acc = append(acc, v)
			}
		}
	}

	acc = append(acc, footerLines...)

	return strings.Join(acc, "\n")
}

func (m model) centerText(rendered string) string {
	lines := strings.Split(rendered, "\n")
	var acc []string

	for _, v := range lines {
		innerAcc := ""
		// m.sizes.width / 2 * 2 - 25% on each side
		for i := 0; i < m.sizes.width/(2*2); i++ {
			innerAcc += " "
		}
		innerAcc += v
		acc = append(acc, innerAcc)
	}

	return strings.Join(acc, "\n")
}

func parseFooter(footer []string) []string {
	var acc []string
	for _, v := range footer {
		result := strings.Split(v, "\n")
		acc = append(acc, result...)

	}

	return acc
}
