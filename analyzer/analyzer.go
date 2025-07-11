// Package analyzer analizes content on the page and extracts
// node thats most likely to contain many <p> tags (article main content)
package analyzer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

type DensityAnalyzer struct {
	page *rod.Page
}

func NewDensityAnalyzer(page *rod.Page) *DensityAnalyzer {
	return &DensityAnalyzer{page: page}
}

type MainArticle struct {
	Content ContentBlock
	Title   string
	Author  string
}

type ContentBlock struct {
	Element     *rod.Element
	TextContent string
	LinkCount   int
	TextLength  int
	Density     float64
	TagName     string
	Area        float64
}

// TODO: Rename
func (da *DensityAnalyzer) AnalyzeContentDensity() (MainArticle, error) {
	// body := da.page.MustElement("body")
	// elems, err := body.Element("section")
	elements, err := da.page.Elements("div, p, article, section, main, aside, header, footer")
	if err != nil {
		return MainArticle{}, err
	}

	var blocks []ContentBlock

	fmt.Println("Calculating elements...")
	for _, element := range elements {
		block, err := da.analyzeElement(element)
		if err != nil {
			continue // Skip problematic elements
		}

		if block.TextLength > 150 {
			blocks = append(blocks, block)
		}
	}

	if len(blocks) == 0 {
		return MainArticle{}, fmt.Errorf("no content blocks found")
	}

	maxDensity := 0.0
	var mainBlock *ContentBlock

	fmt.Println("Weighing score based on elements...")
	for i, block := range blocks {
		score := block.Density

		if block.TagName == "article" || block.TagName == "main" {
			score *= 1.5
		}

		if block.TagName == "option" || block.TagName == "select" || block.TagName == "li" || block.TagName == "ul" {
			score *= 0.01
		}

		// Penalize p tags, likely to be part of article, rather than article itself
		if block.TagName == "p" {
			score *= 0.10
		}

		if block.TagName == "INVALID" {
			score *= 0.05
		}

		// Prefer blocks with substantial text
		if block.TextLength > 1000 {
			score *= 1.2
		}

		if block.TextLength > 5000 {
			score *= 2
		}

		if block.Element.MustEval(`() => this.id`).String() == "auditable" {
			fmt.Println(score)
		}

		blocks[i].Density = score
	}

	fmt.Println("Redistributing score based on parentage...", len(blocks))
	for _, block := range blocks {
		parent, err := block.Element.Parent()
		if err == nil {
			for i, bl := range blocks {
				if da.isSameElement(parent, bl.Element) {
					blocks[i].Density += block.Density * 0.5
				}
				grandPar, err := parent.Parent()
				if err != nil {
					continue
				}
				if da.isSameElement(grandPar, bl.Element) {
					blocks[i].Density += block.Density * 0.25
				}

				grandGrandPar, err := grandPar.Parent()
				if err != nil {
					continue
				}
				if da.isSameElement(grandGrandPar, bl.Element) {
					blocks[i].Density += block.Density * 0.15
				}
			}
		}
	}

	fmt.Println("Choosing the winner...")
	for i, block := range blocks {
		score := block.Density
		if score > maxDensity {
			maxDensity = score
			mainBlock = &blocks[i]
		}
	}

	title := da.getTitle()
	// fmt.Printf("Len %d | Text first 50: %s \t Last 10: %s\n", mainBlock.TextLength, mainBlock.TextContent[:50], mainBlock.TextContent[len(mainBlock.TextContent)-50:])
	// fmt.Printf("Tag name %s | Score: %v | Max density %v \n", mainBlock.TagName, mainBlock.Density, maxDensity)
	art := MainArticle{
		Content: *mainBlock,
		Title:   title,
		Author:  "",
	}

	fmt.Println(mainBlock.TagName, mainBlock.Density, mainBlock.TextContent[0:50], mainBlock.TextContent[len(mainBlock.TextContent)-50:])

	return art, nil
}

func (da *DensityAnalyzer) analyzeElement(element *rod.Element) (ContentBlock, error) {
	textContent := element.MustText()

	cleanText := da.cleanText(textContent)
	cleanLen := len(cleanText)

	// TODO: don't penalize links, if their href is a navigator tag
	// like so <li> <a href="#1" /> </li>
	links := element.MustElements("a")
	linkCount := len(links)

	// for _, v := range links {
	// 	href := v.MustEval(`() => this.href`).String()
	// 	fmt.Println(href)
	// }

	box := element.MustShape().Box()
	var area float64
	if box != nil {
		area = box.Width * box.Height
	}

	tagName := element.MustEval(`() => this.tagName.toLowerCase()`).String()
	if tagName == "" {
		tagName = "INVALID"
	}

	density := da.calculateDensity(cleanLen, linkCount, area)

	block := ContentBlock{
		Element:     element,
		TextContent: cleanText,
		LinkCount:   linkCount,
		TextLength:  cleanLen,
		Area:        area,
		TagName:     tagName,
		Density:     density,
	}

	return block, nil
}

func (da *DensityAnalyzer) cleanText(text string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(strings.TrimSpace(text), " ")

	// Remove common non-content patterns
	patterns := []string{
		`\b(click here|read more|continue reading|share|tweet|like|follow)\b`,
		`\b\d{1,2}[\/\-]\d{1,2}[\/\-]\d{2,4}\b`, // dates
		`\b\w+@\w+\.\w+\b`,                      // emails
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		cleaned = re.ReplaceAllString(cleaned, "")
	}

	return strings.TrimSpace(cleaned)
}

func (da *DensityAnalyzer) calculateDensity(textLength, linkCount int, area float64) float64 {
	if area <= 0 {
		return 0
	}

	baseDensity := float64(textLength) / area

	linkPenalty := 1.0
	if textLength > 0 {
		linkRatio := float64(linkCount) / float64(textLength) * 100
		if linkRatio > 5 { // More than 5% links to text ratio
			linkPenalty = 1.0 - (linkRatio-10)/100
			if linkPenalty < 0.1 {
				linkPenalty = 0.1
			}
		}
	}

	return baseDensity * linkPenalty * 1000
}

// ObjectId will not persist in between calls
func (da *DensityAnalyzer) isSameElement(el1, el2 *rod.Element) bool {
	prop1 := el1.MustEval(`() => 
		this.id + "|" 
		+ this.tagName + '|' 
		+ this.className + '|' 
		+ this.childElementCount + '|'
		+ this.clientHeight + '|'
		+ this.clientWidth + '|'
		+ this.textContent.substring(0, 50)
	`).String()

	prop2 := el2.MustEval(`() => 
		this.id + "|" 
		+ this.tagName + '|' 
		+ this.className + '|' 
		+ this.childElementCount + '|'
		+ this.clientHeight + '|'
		+ this.clientWidth + '|'
		+ this.textContent.substring(0, 50)
	`).String()

	return prop1 == prop2
}

func Run(link string) error {
	browser := rod.New().NoDefaultDevice().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(link)
	page.MustWaitLoad()

	analyzer := NewDensityAnalyzer(page)
	_, err := analyzer.AnalyzeContentDensity()
	if err != nil {
		return fmt.Errorf("error running content analyzer %w", err)
	}
	return nil
}
