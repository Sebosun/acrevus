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
	url, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("error parsing link%w", err)
	}

	browser := rod.New().NoDefaultDevice().MustConnect()
	scrapper := NewScrapper(browser)
	fmt.Println("Host ", url.Host)

	switch url.Host {
	case "substack.com":
		scrapper.substackParser(link)
	default:
		scrapper.generalParserAnalyze(link)
	}

	return nil
}
