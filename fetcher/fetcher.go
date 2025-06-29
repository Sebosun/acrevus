// Package fetcher responsible for fetching articles
package fetcher

import (
	"fmt"
	"sebosun/acrevus-go/storage"
	"strconv"
	"time"

	"github.com/go-rod/rod"
)

func InitFetcher(url string) error {
	browser := rod.New().NoDefaultDevice().MustConnect()
	fmt.Println("Starting parsing: ", url)
	// some random substackpost
	page := browser.MustPage(url).MustWindowFullscreen()

	var title string
	var subtitle string
	var text string

	page.MustElement("article")
	title = page.MustElement("article .post-title").MustText()
	subtitle = page.MustElement("article .subtitle").MustText()

	actualArticle := page.MustElement(".markup")
	elems := actualArticle.MustElements(":not(.subscription-widget-wrap)")

	htmlAcc := []string{}
	for _, v := range elems {
		// nodeType := v.Object.Type
		curHTML, err := v.HTML()
		if err != nil {
			return fmt.Errorf("err getting elements %w", err)
		} else {
			htmlAcc = append(htmlAcc, curHTML)
		}
	}

	fmt.Printf("\n")
	fmt.Printf("Title %s \n", title)
	fmt.Printf("Subtitle %s \n", subtitle)
	fmt.Printf("Text %s \n", trimText(text))

	saveTime := strconv.Itoa(time.Now().YearDay())
	fileName := fmt.Sprintf("%s - %s.html", title, saveTime)

	rawText, err := actualArticle.Text()
	if err != nil {
		return fmt.Errorf("err reading article text %w", err)
	}

	entry := storage.Entry{Title: title, Subtitle: subtitle, Path: fileName, OriginalURL: url, RawText: rawText}

	fmt.Println("Saving...")
	err = storage.SaveArticle(fileName, entry, htmlAcc)
	if err != nil {
		return err
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
