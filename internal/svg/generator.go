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
	baseHeight           = 114.5
	heightStep           = 20.0
	ignoredLanguagesPath = "config/ignored_languages.json"
	templateName         = "template.svg"
	templatePath         = "internal/svg/template.svg"
)

type SVGData struct {
	Theme         Theme
	Height        float64
	Header        string
	Languages     []stats.Lang
	LanguageCount int
}

func GenerateSVG(c *gin.Context) {
	theme := c.DefaultQuery("theme", DefaultTheme)
	header := c.DefaultQuery("header", "Languages")

	languages, err := stats.FetchStats(ignoredLanguagesPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching stats")
		return
	}

	languageCount := len(languages)

	data := SVGData{
		Theme:         GetTheme(theme),
		Height:        calculateSVGHeight(languageCount),
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
	// Height increases by 20 every 2 languages
	steps := (languageCount + 1) / 2
	return baseHeight + float64(steps-1)*heightStep
}

func generateSVG(data SVGData) (string, error) {
	tmpl := template.New(templateName).Funcs(template.FuncMap{
		"sumPrev": sumPreviousPercent,
	})

	tmpl, err := tmpl.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func sumPreviousPercent(languages []stats.Lang, idx int) float64 {
	sum := 0.0
	for i := range idx {
		sum += languages[i].Percent
	}
	return sum
}
