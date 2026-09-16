// Package analyzer analizes content on the page and extracts
// node thats most likely to contain many <p> tags (article main content)
package analyzer

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

const pageLoadTimeout = 30 * time.Second

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
	RawHTML string
}

type ContentBlock struct {
	Element                 *rod.Element
	TextContent             string
	LinkCount               int
	TextLength              int
	Density                 float64
	TagName                 string
	Area                    float64
	CachedFingerprintString string
}

func (da *DensityAnalyzer) ParseContentDensity() (MainArticle, error) {
	elements, err := da.page.Elements("body, div, p, article, section, main, aside, header, footer")
	if err != nil {
		return MainArticle{}, err
	}

	var blocks []ContentBlock

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

	da.weighScoreByTag(&blocks)
	da.redistributeToParents(&blocks)

	for i, block := range blocks {
		score := block.Density
		if score > maxDensity {
			maxDensity = score
			mainBlock = &blocks[i]
		}
	}

	// TODO: extract to separate function
	da.clean(mainBlock.Element, "form")
	da.clean(mainBlock.Element, "fieldset")
	da.clean(mainBlock.Element, "object")
	da.clean(mainBlock.Element, "embed")
	da.clean(mainBlock.Element, "footer")
	da.clean(mainBlock.Element, "link")
	da.clean(mainBlock.Element, "aside")
	da.clean(mainBlock.Element, "iframe")
	da.clean(mainBlock.Element, "input")
	da.clean(mainBlock.Element, "textarea")
	da.clean(mainBlock.Element, "select")
	da.clean(mainBlock.Element, "button")
	da.cleanBr(mainBlock.Element)

	title := da.getTitle()

	rawHTML := mainBlock.Element.MustHTML()
	rawHTML = cleanStyle(cleanClass(rawHTML))

	art := MainArticle{
		Content: *mainBlock,
		Title:   title,
		Author:  "",
		RawHTML: rawHTML,
	}

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
		TextRegex[ClickHere],
		TextRegex[Dates],
		TextRegex[Emails],
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

func SaveTempToDrive(art MainArticle) error {
	file, err := os.Create("./temp.html")
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write([]byte(art.RawHTML))

	return nil
}

func Run(link string) (MainArticle, error) {
	browser := rod.New().NoDefaultDevice().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(link)
	if err := page.Timeout(pageLoadTimeout).WaitLoad(); err != nil {
		return MainArticle{}, fmt.Errorf("wait for page load: %w", err)
	}

	analyzer := NewDensityAnalyzer(page)
	art, err := analyzer.ParseContentDensity()
	fmt.Println("We got out, not we need to get back...")
	if err != nil {
		return MainArticle{}, fmt.Errorf("error running content analyzer %w", err)
	}

	return art, nil
}
