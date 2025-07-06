package storage

import (
	"os"
	"path"
)

func getBasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, ".local", "share", PROJECT_NAME), nil
}

func getEntriesJSONPath() (string, error) {
	filename := "entries.json"
	baseDir, err := getBasePath()
	if err != nil {
		return "", err
	}

	saveDir := path.Join(baseDir, filename)

	return saveDir, nil
}

func getArticlesPath() (string, error) {
	baseDir, err := getBasePath()
	if err != nil {
		return "", err
	}
	articlesDir := path.Join(baseDir, "articles")
	return articlesDir, nil
}
