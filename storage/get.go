package storage

import (
	"encoding/json"
	"os"
)

func GetFileData() (FileData, error) {
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
