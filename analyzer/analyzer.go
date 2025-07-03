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

type ContentBlock struct {
	Element     *rod.Element
	TextContent string
	LinkCount   int
	TextLength  int
	Density     float64
	TagName     string
	Area        float64
}

func (da *DensityAnalyzer) AnalyzeContentDensity() error {
	// body := da.page.MustElement("body")
	// elems, err := body.Element("section")
	elements, err := da.page.Elements("div, p, article, section, main, aside, header, footer")
	if err != nil {
		return err
	}

	var blocks []ContentBlock

	for _, element := range elements {
		block, err := da.analyzeElement(element)
		if err != nil {
			continue // Skip problematic elements
		}

		if block.TextLength > 10 { // Only consider blocks with meaningful text
			blocks = append(blocks, block)
		}
	}

	if len(blocks) == 0 {
		return fmt.Errorf("no content blocks found")
	}

	maxDensity := 0.0
	var mainBlock *ContentBlock

	for i, block := range blocks {
		// Boost score for semantic tags
		score := block.Density
		if block.TagName == "article" || block.TagName == "main" {
			score *= 1.5
		}

		if block.TagName == "option" || block.TagName == "select" || block.TagName == "li" || block.TagName == "ul" {
			score *= 0.01
		}

		if block.TagName == "INVALID" {
			score *= 0.05
		}

		// Prefer blocks with substantial text
		if block.TextLength > 200 {
			score *= 1.2
		}

		if score > maxDensity {
			maxDensity = score
			mainBlock = &blocks[i]
		}
	}

	fmt.Printf("Len %d | Text first 10: %s\n", mainBlock.TextLength, mainBlock.TextContent)
	fmt.Printf("Tag name %s | Score: %v | Max density %v \n", mainBlock.TagName, mainBlock.Density, maxDensity)

	return nil
}

func (da *DensityAnalyzer) analyzeElement(element *rod.Element) (ContentBlock, error) {
	textContent, err := element.Text()
	if err != nil {
		return ContentBlock{}, err
	}

	cleanText := da.cleanText(textContent)
	cleanLen := len(cleanText)

	links, _ := element.Elements("a")
	linkCount := len(links)

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

	return baseDensity * linkPenalty
}

func Run(link string) error {
	browser := rod.New().NoDefaultDevice().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(link)
	page.MustWaitLoad()

	page.MustScreenshot("a.png")
	analyzer := NewDensityAnalyzer(page)
	err := analyzer.AnalyzeContentDensity()
	if err != nil {
		return fmt.Errorf("error running content analyzer %w", err)
	}
	return nil
}
