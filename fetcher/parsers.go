package fetcher

import (
	"fmt"
	"sebosun/acrevus-go/storage"
	"strconv"
	"time"
)

func (br *Scapper) generalParser(link string) error {
	fmt.Println("Starting general parsing: ", link)
	page := br.browser.MustPage(link).MustWindowFullscreen()
	page.MustScreenshot("general-parser.png")

	var title string
	var subtitle string
	var text string

	text = page.MustElement("article").MustText()
	title = page.MustElement("article h1").MustText()
	subtitle = page.MustElement("article > h2,h3,h4,h5").MustText()

	fmt.Printf("Title %s \n", title)
	fmt.Printf("Subtitle %s \n", subtitle)
	fmt.Printf("Text %s \n", trimText(text))

	htmlAcc := []string{}
	for _, v := range page.MustElements("article") {
		// nodeType := v.Object.Type
		curHTML, err := v.HTML()
		if err != nil {
			return fmt.Errorf("err getting elements %w", err)
		} else {
			htmlAcc = append(htmlAcc, curHTML)
		}
	}

	saveTime := strconv.Itoa(time.Now().YearDay())
	fileName := fmt.Sprintf("%s - %s.html", title, saveTime)

	entry := storage.Entry{Title: title, Subtitle: subtitle, Path: fileName, OriginalURL: link, RawText: text}

	err := storage.SaveArticle(fileName, entry, htmlAcc)
	if err != nil {
		return err
	}

	return nil
}

func (br *Scapper) substackParser(link string) error {
	fmt.Println("Starting substack parsing: ", link)

	page := br.browser.MustPage(link).MustWindowFullscreen()
	page.MustScreenshot("substack-parser.png")

	var title string
	var subtitle string
	var text string

	page.MustElement("article")
	title = page.MustElement("article .post-title").MustText()
	subtitle = page.MustElement("article .subtitle").MustText()

	actualArticle := page.MustElement(".markup")
	elems := actualArticle.MustElements(":not(.subscription-widget-wrap)")

	fmt.Println("Going through nodes", link)
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

	entry := storage.Entry{Title: title, Subtitle: subtitle, Path: fileName, OriginalURL: link, RawText: rawText}

	fmt.Println("Saving...")
	err = storage.SaveArticle(fileName, entry, htmlAcc)
	if err != nil {
		return err
	}

	return nil
}
