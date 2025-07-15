package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func SaveArticle(filename string, entry Entry, data []string) error {
	fileData, err := GetFileData()
	fileData.Entries = append(fileData.Entries, entry)
	if err != nil {
		return err
	}
	overwrite(fileData)
	saveArticle(filename, data)

	fmt.Printf("Successfully wrote %d strings to %s\n", len(data), filename)
	return nil
}

func saveArticle(filename string, html []string) error {
	articlesPath, err := GetArticlesPath()
	if err != nil {
		return err
	}

	newArtPath := path.Join(articlesPath, filename)
	articleFile, err := os.Create(newArtPath)
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
	for _, str := range html {
		_, err := articleFile.WriteString(str + "\n")
		if err != nil {
			return fmt.Errorf("failed to write string: %w", err)
		}
	}
	articleFile.WriteString("</div>" + "\n")
	return nil
}

func overwrite(data FileData) error {
	jsonPath, err := GetEntriesJSONPath()
	if err != nil {
		return err
	}

	jsonFile, err := os.Create(jsonPath)
	if err != nil {
		return err
	}

	defer jsonFile.Close()
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = jsonFile.Write(result)
	if err != nil {
		return err
	}

	return nil
}
