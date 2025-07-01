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
	xd, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("error parsing link%w", err)
	}

	browser := rod.New().NoDefaultDevice().MustConnect()
	scrapper := NewScrapper(browser)
	fmt.Println("Host ", xd.Host)

	switch xd.Host {
	case "substack.com":
		scrapper.substackParser(link)
	default:
		scrapper.generalParser(link)
	}

	return nil
}

func trimText(input string) string {
	if len(input) == 0 {
		return ""
	}

	acc := ""
	for i := 0; i < 50 || i < len(input); i++ {
		acc += string(input[i])
	}
	return acc
}
