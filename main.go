package main

import (
	"log"
)

func main() {
	_, err := ensureFilesExist()
	if err != nil {
		log.Fatalf("Failed to get stored data: %v", err)
	}
	err = initFetcher(true)
	if err != nil {
		log.Fatalf("Error fetching %v", err)
	}
}
