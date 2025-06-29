package main

import (
	"log"
	"sebosun/acrevus-go/cmd"
	"sebosun/acrevus-go/storage"
)

func main() {
	_, err := storage.EnsureFilesExist()
	if err != nil {
		log.Fatalf("Failed to get stored data: %v", err)
	}

	cmd.Execute()
}
