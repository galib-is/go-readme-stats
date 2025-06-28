package api

import (
	"go-readme-stats/internal/stats"
	"go-readme-stats/internal/svg"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ignoredLanguagesPath = "config/ignored_languages.json"

func GetLanguageStats(c *gin.Context) {
	theme := c.DefaultQuery("theme", svg.DefaultTheme)
	header := c.DefaultQuery("header", "Languages")

	languages, err := stats.FetchStats(ignoredLanguagesPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching stats")
		return
	}

	svgContent, err := svg.Generate(theme, header, languages)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating SVG")
		return
	}

	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svgContent)
}
