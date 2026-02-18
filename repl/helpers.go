package repl

import "strings"

func (m model) renderOnlyVisible(article string, footer []string) string {
	lines := strings.Split(article, "\n")
	footerLines := parseFooter(footer)

	var acc []string
	footLen := len(footerLines)

	for i, v := range lines {
		isBeforeView := i + m.sizes.height >= m.cursor
		isAfterView := i < m.sizes.height+m.cursor-footLen

		if isBeforeView && isAfterView {
			if i == m.cursor {
				// curRead := currentlyReadingStyle.Render(v[1:])
				// acc = append(acc, ">"+curRead)
				acc = append(acc, ">"+v)
			} else {
				acc = append(acc, v)
			}
		}
	}

	acc = append(acc, footerLines...)

	return strings.Join(acc, "\n")
}

func (m model) centerText(rendered string, paddingProcent float64) string {
	lines := strings.Split(rendered, "\n")
	var acc []string

	linesWidthProcent := (1 - paddingProcent) / 2
	maxWidth := int(float32(m.sizes.width) * float32(linesWidthProcent))

	for _, v := range lines {
		var innerAcc strings.Builder
		// m.sizes.width / 2 * 2 - 25% of entire thing
		for range maxWidth {
			innerAcc.WriteString(" ")
		}
		innerAcc.WriteString(v)
		acc = append(acc, innerAcc.String())
	}

	return strings.Join(acc, "\n")
}

func (m *model) updateViewHeight(height int) {
	m.viewHeight = height
}

func parseFooter(footer []string) []string {
	var acc []string
	for _, v := range footer {
		result := strings.Split(v, "\n")
		acc = append(acc, result...)

	}

	return acc
}
