// Package api provides HTTP handlers

package api

import (
	"go-readme-stats/internal/stats"
	"go-readme-stats/internal/svg"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ignoredLanguagesPath = "config/ignored_languages.json"

// Add these variables for dependency injection and testing
var (
	FetchStats  = stats.FetchStats
	GenerateSVG = svg.Generate
)

// GetLanguageStats handles GET requests for SVG generation.
// Query parameters:
//   - theme: SVG theme name (defaults to DefaultTheme - dark)
//   - header: Custom header text (defaults to "Languages")
//
// Returns an SVG image with HTTP 200 on success, or HTTP 500 on error.
func GetLanguageStats(c *gin.Context) {
	theme := c.DefaultQuery("theme", svg.DefaultTheme)
	header := c.DefaultQuery("header", "Languages")

	languages, err := FetchStats(ignoredLanguagesPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching stats")
		return
	}

	svgContent, err := GenerateSVG(theme, header, languages)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating SVG")
		return
	}

	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svgContent)
}
