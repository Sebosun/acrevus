package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-rod/rod"
)

func InitFetcher(url string, run bool) error {
	if !run {
		return nil
	}

	browser := rod.New().NoDefaultDevice().MustConnect()
	// some random substackpost
	page := browser.MustPage("https://substack.com/home/post/p-166878772").MustWindowFullscreen()

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
		nodeType := v.Object.Type
		text, err := v.HTML()
		if err != nil {
			log.Fatalf("err getting elements %v", err)
		} else {
			fmt.Println("Type ", nodeType)
			htmlAcc = append(htmlAcc, text)
		}
	}

	fmt.Printf("\n")
	fmt.Printf("Title %s \n", title)
	fmt.Printf("Subtitle %s \n", subtitle)
	fmt.Printf("Text %s \n", trimText(text))

	saveTime := strconv.Itoa(time.Now().YearDay())
	fileName := fmt.Sprintf("%s - %s.html", title, saveTime)

	entry := Entry{Title: title, Subtitle: subtitle, Path: fileName, OriginalURL: url}

	err := saveArticle(fileName, entry, htmlAcc)
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
