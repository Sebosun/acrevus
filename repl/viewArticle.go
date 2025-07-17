package repl

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/net/html"
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

	result, err := ViewHTMLArticle(m.articleRawHtml)
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
	visible := m.renderOnlyVisible(stylized, footer)

	file, err := os.Create("./temp-render.html")
	if err != nil {
		return ""
	}
	defer file.Close()
	file.Write([]byte(s))

	return visible
}

func ViewHTMLArticle(rawHTML string) (DisplayNode, error) {
	n, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return DisplayNode{}, err
	}

	nodes := runTextParser(n, 0)

	return nodes, nil
}

func getTexts(n DisplayNode) string {
	s := ""
	re := regexp.MustCompile(`w+`)

	switch n.ParentNode {
	case atom.Span:
		attachString := strings.TrimLeft(n.TextContent, "\t")
		attachString = strings.TrimLeft(attachString, " ")
		attachString = strings.TrimLeft(attachString, "\n")
		hasWords := re.Match([]byte(attachString))
		if hasWords {
			s += attachString + " "
		}

	case atom.P:
		s += "\n"
		s += n.TextContent
		s += "\n"
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5:
		s += "\n"
		s += headerStyle.Render(n.TextContent)
		s += "\n"
	case atom.Li:
		style := lipgloss.NewStyle().
			Bold(false).
			Padding(0, 2).
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#7D56F4"))
		s += style.Render(n.TextContent) + " - LI"
	default:
		re := regexp.MustCompile(`[a-zA-Z0-9]`)
		result := re.Match([]byte(n.TextContent))
		if result {
			s += n.TextContent + " "
		} else {

		}
	}

	for _, child := range n.Children {
		txt := getTexts(child)
		s += txt
	}

	return s
}
