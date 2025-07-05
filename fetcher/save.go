package fetcher

import (
	"fmt"
	"sebosun/acrevus-go/storage"
	"strconv"
	"time"
)

type SaveData struct {
	title    string
	subtitle string
	url      string
	text     string
}

func saveToDrive(data SaveData, html []string) error {
	saveTime := strconv.Itoa(time.Now().YearDay())
	fileName := fmt.Sprintf("%s - %s.html", data.title, saveTime)

	entry := storage.Entry{Title: data.title, Subtitle: data.subtitle, Path: fileName, OriginalURL: data.url, RawText: data.text}
	err := storage.SaveArticle(fileName, entry, html)
	if err != nil {
		return err
	}
	return nil
}
