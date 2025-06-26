package svg

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"go-readme-stats/internal/stats"

	"github.com/gin-gonic/gin"
)

const (
	baseHeight   = 114.5
	mediumHeight = 134.5
	largeHeight  = 154.5
)

type SVGData struct {
	Height        float64
	Header        string
	Languages     []stats.Lang
	LanguageCount int
}

func GenerateSVG(c *gin.Context) {
	header := "Languages"
	username := "galib-i"
	ignoredLangsPath := "ignored_languages.json"

	languages, err := stats.FetchStats(username, ignoredLangsPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching stats")
		return
	}

	languageCount := len(languages)
	height := calculateSVGHeight(languageCount)

	data := SVGData{
		Height:        height,
		Header:        header,
		Languages:     languages,
		LanguageCount: languageCount,
	}

	svg, err := generateSVG(data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating SVG")
		return
	}

	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svg)
}

func calculateSVGHeight(languageCount int) float64 {
	switch {
	case languageCount <= 2:
		return baseHeight
	case languageCount <= 4:
		return mediumHeight
	default:
		return largeHeight
	}
}

func generateSVG(data SVGData) (string, error) {
	tmpl := template.New("template.svg").Funcs(template.FuncMap{
		"sumPrev": sumPreviousPercent,
	})

	tmpl, err := tmpl.ParseFiles("internal/svg/template.svg")
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func sumPreviousPercent(langs []stats.Lang, idx int) float64 {
	sum := 0.0
	for i := range idx {
		sum += langs[i].Percent
	}
	return sum
}
