package main

import (
	"fmt"
	"log"
	"os"

	"go-readme-stats/scripts"
)

func main() {
	output := "internal/data/colours.json"
	if _, err := os.Stat(output); os.IsNotExist(err) {
		err := scripts.FetchLanguageColours(output)
		if err != nil {
			log.Fatalf("failed to fetch and convert language colours: %v", err)
		}
		log.Printf("Conversion complete. Saved to %s", output)
	} else {
		fmt.Println("file exists")
	}
}
