package storage

import (
	"os"
	"path"
)

func GetBasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, ".local", "share", PROJECT_NAME), nil
}

func GetEntriesJSONPath() (string, error) {
	filename := "entries.json"
	baseDir, err := GetBasePath()
	if err != nil {
		return "", err
	}

	saveDir := path.Join(baseDir, filename)

	return saveDir, nil
}

func GetArticlesPath() (string, error) {
	baseDir, err := GetBasePath()
	if err != nil {
		return "", err
	}
	articlesDir := path.Join(baseDir, "articles")
	return articlesDir, nil
}
