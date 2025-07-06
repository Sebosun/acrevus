package storage

import (
	"errors"
	"os"
)

func EnsureFilesExist() (FileData, error) {
	err := ensureFolderExist()
	if err != nil {
		return FileData{}, err
	}

	created, err := ensureFileExists()
	if err != nil && !created {
		return FileData{}, err
	}

	data, err := GetFileData()
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
