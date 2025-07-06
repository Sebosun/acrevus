package storage

import (
	"fmt"
	"os"
	"path"
)

func DeleteData(entry Entry) (FileData, error) {
	data, err := GetFileData()
	if err != nil {
		return FileData{}, fmt.Errorf("delete data entry %v", err)
	}

	deleteArticle(entry.Path)

	newData := FileData{}
	for _, v := range data.Entries {
		if v.OriginalURL == entry.OriginalURL {
			continue
		}

		newData.Entries = append(newData.Entries, v)
	}

	err = overwrite(newData)
	if err != nil {
		return FileData{}, fmt.Errorf("error deleting data while overwrtigin %v", err)
	}

	return newData, nil
}

func deleteArticle(relPath string) error {
	articlesPath, err := getArticlesPath()
	if err != nil {
		return fmt.Errorf("delete data get artc error %v", err)
	}

	htmlFilePath := path.Join(articlesPath, relPath)
	err = os.Remove(htmlFilePath)
	if err != nil {
		return fmt.Errorf("deleting html file error %v", err)
	}
	return nil
}
