package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

type Entry struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Path     string `json:"path"`
}

type FileData struct {
	Entries []Entry `json:"entries"`
}

type DB struct {
	fileData FileData
}

const PROJECT_NAME = "acrevus"

func ensureFilesExist() (FileData, error) {
	err := ensureFolderExist()
	if err != nil {
		return FileData{}, err
	}

	created, err := ensureFileExists()
	if err != nil && !created {
		return FileData{}, err
	}

	data, err := getFileData()
	if err != nil {
		return FileData{}, err
	}
	return data, nil
}

func ensureFileExists() (bool, error) {
	path, err := getEntriesJSONPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		_, err := createDataFile(path)
		if err != nil {
			return true, err
		}
		return false, err
	} else if err != nil {
		return false, err
	}

	return false, nil
}

func ensureFolderExist() error {
	articlesPath, err := getArticlesPath()
	if err != nil {
		return nil
	}

	info, err := os.Stat(articlesPath)
	if errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(articlesPath, 0755)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("articles folder is not a folder")
	}

	return nil
}

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

func getFileData() (FileData, error) {
	path, err := getEntriesJSONPath()
	if err != nil {
		return FileData{}, err
	}

	var fileData FileData
	res, err := os.ReadFile(path)
	if err != nil {
		return FileData{}, err
	}

	err = json.Unmarshal(res, &fileData)
	if err != nil {
		return FileData{}, err
	}

	return fileData, nil
}

func saveArticle(filename string, entry Entry, data []string) error {
	fileData, err := getFileData()
	fileData.Entries = append(fileData.Entries, entry)
	if err != nil {
		return err
	}

	getBasePath, err := getBasePath()
	if err != nil {
		return err
	}

	jsonPath, err := getEntriesJSONPath()
	if err != nil {
		return err
	}

	jsonFile, err := os.Create(jsonPath)
	if err != nil {
		return err
	}

	defer jsonFile.Close()
	result, err := json.Marshal(fileData)
	if err != nil {
		return err
	}
	_, err = jsonFile.Write(result)
	if err != nil {
		return err
	}

	saveDir := path.Join(getBasePath, "articles", filename)
	articleFile, err := os.Create(saveDir)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		if closeErr := articleFile.Close(); closeErr != nil {
			fmt.Printf("Error closing file: %v\n", closeErr)
		}
	}()

	// wrapping the entire thing in div just in case
	articleFile.WriteString("<div>" + "\n")
	for _, str := range data {
		_, err := articleFile.WriteString(str + "\n")
		if err != nil {
			return fmt.Errorf("failed to write string: %w", err)
		}
	}
	articleFile.WriteString("</div>" + "\n")

	fmt.Printf("Successfully wrote %d strings to %s\n", len(data), filename)
	return nil
}
