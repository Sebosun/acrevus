// Package analyzer analizes content on the page and extracts
// node thats most likely to contain many <p> tags (article main content)
package analyzer

import (
	"fmt"

	"github.com/go-rod/rod"
)

type DensityAnalyzer struct {
	page *rod.Page
}

func NewDensityAnalyzer(page *rod.Page) *DensityAnalyzer {
	return &DensityAnalyzer{page: page}
}

func (da *DensityAnalyzer) AnalyzeContentDensity() error {
	elements := da.page.MustElement("div").MustText()

	fmt.Printf("%d \n", len(elements))
	// for i, v := range elements {
	// 	fmt.Println(i, v.MustText())
	// }
	return nil
}

func Run(link string) error {
	browser := rod.New().NoDefaultDevice().MustConnect()
	page := browser.MustPage(link).MustWindowFullscreen()

	page.MustScreenshot("a.png")
	analyzer := NewDensityAnalyzer(page)
	err := analyzer.AnalyzeContentDensity()
	if err != nil {
		return fmt.Errorf("error running content analyzer %w", err)
	}
	return nil
}
