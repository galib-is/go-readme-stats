// Package main starts a web server to generate SVGs for the user's GitHub language statistics.

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go-readme-stats/internal/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	router := gin.Default()
	router.GET("/langs", api.GetLanguageStats)
	router.Run("localhost:8080")
}
