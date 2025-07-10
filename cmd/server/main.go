// Package main starts a web server to generate SVGs for the user's GitHub language statistics.

package main

import (
	"go-readme-stats/internal/api"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	/*
		if err := scripts.FetchLanguageColours(); err != nil {
			log.Fatalf("Failed to fetch language colours: %v", err)
		} else {
			log.Println("Successfully fetched language colours.")
		}
	*/

	router := gin.Default()
	router.GET("/langs", api.GetLanguageStats)
	router.Run("localhost:8080")
}
