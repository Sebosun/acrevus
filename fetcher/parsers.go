package fetcher

import (
	"fmt"
	"time"

	"sebosun/acrevus-go/analyzer"
)

const pageLoadTimeout = 30 * time.Second

func (br *Scapper) generalParserAnalyze(link string) error {
	page := br.browser.MustPage(link)
	if err := page.Timeout(pageLoadTimeout).WaitLoad(); err != nil {
		return fmt.Errorf("wait for page load: %w", err)
	}
	densityAnalyzer := analyzer.NewDensityAnalyzer(page)
	article, err := densityAnalyzer.ParseContentDensity()

	if err != nil {
		return fmt.Errorf("error running content analyzer %w", err)
	}

	data := SaveData{
		title:    article.Title,
		subtitle: "",
		url:      link,
		text:     article.Content.TextContent,
	}

	err = saveToDrive(data, []string{article.RawHTML})
	if err != nil {
		return fmt.Errorf("error saving content via analyzer %w", err)
	}

	return nil
}
