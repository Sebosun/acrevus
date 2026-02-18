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
	purple     = lipgloss.Color("99")
	purpleLink = lipgloss.Color("120")
	gray       = lipgloss.Color("245")
	lightGray  = lipgloss.Color("241")

	linkStyle             = lipgloss.NewStyle().Foreground(purpleLink).Bold(true).Align(lipgloss.Center)
	headerStyle           = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	currentlyReadingStyle = lipgloss.NewStyle().Background(lightGray).Align(lipgloss.Center)
)

func (m model) ViewArticle() string {
	s := ""

	// art := m.articles[m.selected]
	// s += m.articleRawHtml

	// This likely should happen when the article is selected, not each time article is being viewed
	art := m.articles[m.selected]
	if len(art.Title) > 0 {
		s += headerStyle.Render(art.Title) + "\n \n"
	}

	if len(art.Subtitle) > 0 {
		s += headerStyle.Render(art.Subtitle) + "\n \n"
	}

	result, err := getParsedArticles(m.articleRawHTML)
	if err != nil {
		return "Error parsing the article"
	}

	s += getTexts(result)

	// The footer
	var footer []string
	//
	// footer = append(footer, fmt.Sprintf("\t \t \t  Width: %d, Height: %d \t     Press q to quit", m.sizes.width, m.sizes.height))
	footer = append(footer, "Use hjkl to move. U - scroll up, D - scroll down. q - quit")
	footer = append(footer, fmt.Sprintf("c:[%d] h:[%d] S:[%d]", m.cursor, m.sizes.height, m.selected))

	paddingProcent := 0.80

	lipglossWidth := m.sizes.width
	if m.opts.isCentered {
		lipglossWidth = int(float32(m.sizes.width) * float32(paddingProcent))
	}

	style := lipgloss.NewStyle().
		Bold(false).
		Width(lipglossWidth).
		Padding(1, 2)
	stylized := style.Render(s)
	visible := m.renderOnlyVisible(stylized, footer)

	if m.opts.isCentered {
		visible = m.centerText(visible, paddingProcent)
	}

	if m.opts.debug {
		file, err := os.Create("./temp-render.html")
		if err != nil {
			return ""
		}
		defer file.Close()
		file.Write([]byte(s))
	}

	return visible
}

func getTexts(n DisplayNode) string {
	s := ""
	re := regexp.MustCompile(`[a-zA-Z0-9]`)

	var child_text strings.Builder
	for _, child := range n.Children {
		txt := getTexts(child)
		child_text.WriteString(txt)
	}

	text := n.TextContent + child_text.String()
	switch n.NodeType {
	case atom.Span:
		attachString := strings.TrimLeft(text, "\t")
		attachString = strings.TrimLeft(text, " ")
		attachString = strings.TrimLeft(text, "\n")
		hasWords := re.Match([]byte(text))
		if hasWords {
			s += attachString
		}
	case atom.P:
		s += text
		s += "\n \n"
	case atom.A:
		s += linkStyle.Render(text)
	case atom.Button:
		s += headerStyle.Render(text)
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
