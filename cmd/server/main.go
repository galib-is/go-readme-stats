package main

import (
	"log"

	"github.com/joho/godotenv"

	"go-readme-stats/internal/svg"
	"go-readme-stats/scripts"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	output := "internal/data/colours.json"

	if err := scripts.EnsureLanguageColours(output); err != nil {
		log.Fatalf("could not ensure colours file: %v", err)
	}

	router := gin.Default()
	router.GET("/langs", svg.GenerateSVG)
	router.Run("localhost:8080")
}
