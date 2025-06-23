package svg

import (
	"bytes"
	"net/http"
	"text/template"

	"go-readme-stats/internal/stats"

	"github.com/gin-gonic/gin"
)

func GenerateSVG(c *gin.Context) {
	username := "galib-i"
	ignoredLangsPath := "ignored_languages.json"

	languages := stats.FetchStats(username, ignoredLangsPath)
	svg := generateSVG(languages)
	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svg)
}

func generateSVG(langs []stats.Lang) string {
	tmpl, err := template.ParseFiles("internal/svg/template.svg")
	if err != nil {
		return "<svg><!-- template error --></svg>"
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, langs)
	if err != nil {
		return "<svg><!-- execute error --></svg>"
	}

	return buf.String()
}
