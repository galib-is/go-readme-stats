package main

import (
	"log"

	"go-readme-stats/internal/svg"
	"go-readme-stats/scripts"

	"github.com/gin-gonic/gin"
)

func main() {
	output := "internal/data/colours.json"

	if err := scripts.EnsureLanguageColours(output); err != nil {
		log.Fatalf("could not ensure colours file: %v", err)
	}

	router := gin.Default()
	router.GET("/langs", svg.GenerateSVG)
	router.Run("localhost:8080")
}
