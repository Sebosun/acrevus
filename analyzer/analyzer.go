// Package analyzer analizes content on the page and extracts
// node thats most likely to contain many <p> tags (article main content)
package analyzer

import (
	"fmt"
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
	RawHTML string
	Metadata Metadata
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
	elements, err := da.page.Elements("div, p, article, section, main, aside, header, footer, h1, h2, h3, h4, h5")
	if err != nil {
		return MainArticle{}, err
	}

	var blocks []ContentBlock

	for _, element := range elements {
		block, err := da.analyzeElement(element)
		if err != nil {
			continue // Skip problematic elements
		}

		blocks = append(blocks, block)
	}

	if len(blocks) == 0 {
		return MainArticle{}, fmt.Errorf("no content blocks found")
	}

	da.weighScoreByTag(&blocks)
	da.redistributeToParents(&blocks)
	mainBlock := getHighestDensityBlock(&blocks)

	da.bulkClean(mainBlock.Element)

	rawHTML := mainBlock.Element.MustHTML()
	rawHTML = cleanStyle(cleanClass(rawHTML))

	metadata, err := GetMetadata(da.page.MustHTML())

	if err != nil {
		return MainArticle{}, err
	}

	art := MainArticle{
		Content: *mainBlock,
		RawHTML: rawHTML,
		Metadata: metadata,
	}

	return art, nil
}

func getHighestDensityBlock(blocks *[]ContentBlock) *ContentBlock {
	maxDensity := 0.0
	var mainBlock *ContentBlock
	for i, block := range *blocks {
		score := block.Density
		if score > maxDensity {
			maxDensity = score
			mainBlock = &(*blocks)[i]
		}
	}
	return mainBlock
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

func (da *DensityAnalyzer) calculateDensity(textLength, linkCount int, area float64) float64 {
	if area <= 0 {
		return 0
	}

	baseDensity := float64(textLength) / area
	linkPenalty := 1.0

	// Safeguard ig, might break some posts that have tons of links to other sources
	if linkCount > 100 {
		return 0.001
	}

	if textLength > 0 {
		linkRatio := float64(linkCount) / float64(textLength) * 100
		if linkRatio > 5.0 { // More than 5% links to text ratio
			linkPenalty = 1.0 - (linkRatio-10)/100

			if linkPenalty < 0.1 {
				linkPenalty = 0.1
			}
		}
	}

	return baseDensity * linkPenalty * 1000
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
	if err != nil {
		return MainArticle{}, fmt.Errorf("error running content analyzer %w", err)
	}

	return art, nil
}
