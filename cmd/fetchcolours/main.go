package main

import (
	"fmt"
	"log"

	"go-readme-stats/scripts"
)

func main() {
	if err := scripts.FetchLanguageColours(); err != nil {
		log.Fatalf("Failed to fetch language colours: %v", err)
	}
	fmt.Println("Successfully fetched language colours.")
}
