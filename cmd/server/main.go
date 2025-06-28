package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go-readme-stats/internal/svg"
	"go-readme-stats/scripts"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	outputPath := "internal/data/colours.json"

	if err := scripts.EnsureLanguageColours(outputPath); err != nil {
		log.Fatalf("could not ensure colours file: %v", err)
	}

	router := gin.Default()
	router.GET("/langs", svg.GenerateSVG)
	router.Run("localhost:8080")
}
