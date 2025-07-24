package repl

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/net/html/atom"
)

var (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")

	headerStyle           = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	currentlyReadingStyle = lipgloss.NewStyle().Background(lightGray).Align(lipgloss.Center)
)

func (m model) ViewArticle() string {
	s := ""

	// art := m.articles[m.selected]
	// s += m.articleRawHtml

	result, err := getParsedArticles(m.articleRawHtml)
	if err != nil {
		return "Error parsing the article"
	}

	s += getTexts(result)

	// The footer
	var footer []string
	//
	footer = append(footer, fmt.Sprintf("\t \t \t  Width: %d, Height: %d \t     Press q to quit", m.sizes.width, m.sizes.height))
	footer = append(footer, fmt.Sprintf("c:[%d] h:[%d] S:[%d]", m.cursor, m.sizes.height, m.selected))
	//
	style := lipgloss.NewStyle().
		Bold(false).
		Width(int(m.sizes.width/2)).
		Padding(1, 2)
	stylized := style.Render(s)
	visible := m.centerText(m.renderOnlyVisible(stylized, footer))

	file, err := os.Create("./temp-render.html")
	if err != nil {
		return ""
	}
	defer file.Close()
	file.Write([]byte(s))

	return visible
}

func getTexts(n DisplayNode) string {
	s := ""
	re := regexp.MustCompile(`[a-zA-Z0-9]`)

	child_text := ""
	for _, child := range n.Children {
		txt := getTexts(child)
		child_text += txt
	}

	text := n.TextContent + child_text
	switch n.NodeType {
	case atom.Span:
		attachString := strings.TrimLeft(text, "\t")
		attachString = strings.TrimLeft(text, " ")
		attachString = strings.TrimLeft(text, "\n")
		hasWords := re.Match([]byte(text))
		if hasWords {
			s += attachString + " "
		}
	case atom.P:
		s += text
		s += "\n \n"
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5:
		s += headerStyle.Render(text)
		s += "\n \n"
	case atom.Li:
		s += "  - " + strings.TrimLeft(text, "\n \n")
	default:
		result := re.Match([]byte(text))
		if result {
			s += text + " "
		} else {
		}
	}

	return cleanText(s)
}

func cleanText(input string) string {
	reNewSpaces := regexp.MustCompile(`\n{3}`)
	reSpaces := regexp.MustCompile(` {2,}`)
	reTwoNew := regexp.MustCompile(`\n `)

	stripped := reTwoNew.ReplaceAll([]byte(input), []byte("\n"))
	stripped = reNewSpaces.ReplaceAll([]byte(stripped), []byte("\n"))
	stripped = reSpaces.ReplaceAll([]byte(stripped), []byte(" "))

	return string(stripped)

}
