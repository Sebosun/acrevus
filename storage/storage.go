// Package storage handles saivng and reading from drive
package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type Entry struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Path        string `json:"path"`
	OriginalURL string `json:"originalUrl"`
	RawText     string `json:"rawText"`
}

type FileData struct {
	Entries []Entry `json:"entries"`
}

type DB struct {
	fileData FileData
}

const PROJECT_NAME = "acrevus"

func IsURLSaved(url string) (bool, error) {
	data, err := GetFileData()
	if err != nil {
		return false, fmt.Errorf("failed to create file: %w", err)
	}

	for _, v := range data.Entries {
		if v.OriginalURL == url {
			return true, nil
		}
	}

	return false, nil
}

func createDataFile(path string) (bool, error) {
	file, err := os.Create(path)
	if err != nil {
		return false, err
	}

	defer file.Close()
	result, err := json.Marshal(FileData{})
	if err != nil {
		return false, err
	}

	_, err = file.Write(result)
	if err != nil {
		return false, err
	}

	return true, nil
}
