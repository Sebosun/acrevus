package repl

import (
	"sebosun/acrevus-go/storage"
)

func (m model) deleteArticle() ([]storage.Entry, error) {
	artToDelet := m.articles[m.selected]
	newFileData, err := storage.DeleteData(artToDelet)
	if err != nil {
		return []storage.Entry{}, err
	}
	return newFileData.Entries, nil
}
