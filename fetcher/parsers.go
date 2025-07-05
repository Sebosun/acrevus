package fetcher

import (
	"fmt"
	"sebosun/acrevus-go/analyzer"
)

func (br *Scapper) generalParserAnalyze(link string) error {
	page := br.browser.MustPage(link)
	page.MustWaitLoad()
	densityAnalyzer := analyzer.NewDensityAnalyzer(page)
	article, err := densityAnalyzer.AnalyzeContentDensity()
	if err != nil {
		return fmt.Errorf("error running content analyzer %w", err)
	}

	data := SaveData{
		title:    article.Title,
		subtitle: "",
		url:      link,
		text:     article.Content.TextContent,
	}

	err = saveToDrive(data, []string{"2", "3"})
	if err != nil {
		return fmt.Errorf("error saving content via analyzer %w", err)
	}

	return nil
}

func (br *Scapper) generalParser(link string) error {
	page := br.browser.MustPage(link).MustWindowFullscreen()
	page.MustScreenshot("general-parser.png")

	text := page.MustElement("article").MustText()
	title := page.MustElement("article h1").MustText()
	subtitle := page.MustElement("article > h2,h3,h4,h5").MustText()

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

	data := SaveData{
		title:    title,
		subtitle: subtitle,
		url:      link,
		text:     text,
	}

	err := saveToDrive(data, htmlAcc)
	if err != nil {
		return fmt.Errorf("error saving general parser article %w", err)
	}

	return nil
}

func (br *Scapper) substackParser(link string) error {
	fmt.Println("Starting substack parsing: ", link)

	page := br.browser.MustPage(link).MustWindowFullscreen()
	page.MustScreenshot("substack-parser.png")

	page.MustElement("article")
	title := page.MustElement("article .post-title").MustText()
	subtitle := page.MustElement("article .subtitle").MustText()

	actualArticle := page.MustElement(".markup")
	text := actualArticle.MustText()

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

	data := SaveData{
		title:    title,
		subtitle: subtitle,
		url:      link,
		text:     text,
	}

	err := saveToDrive(data, htmlAcc)
	if err != nil {
		return fmt.Errorf("Error saving substack parser article %w", err)
	}

	return nil
}
