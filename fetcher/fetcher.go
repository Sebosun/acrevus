// Package fetcher responsible for fetching articles
package fetcher

import (
	"fmt"
	"net/url"

	"github.com/go-rod/rod"
)

type Scapper struct {
	browser *rod.Browser
}

func NewScrapper(browser *rod.Browser) *Scapper {
	return &Scapper{browser: browser}
}

func InitFetcher(link string) error {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("error parsing link%w", err)
	}

	browser := rod.New().NoDefaultDevice().MustConnect()
	scrapper := NewScrapper(browser)

	scrapper.generalParserAnalyze(link)

	return nil
}
