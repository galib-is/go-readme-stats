package svg

import (
	"bytes"
	"net/http"
	"text/template"

	"go-readme-stats/internal/stats"

	"github.com/gin-gonic/gin"
)

type SVGData struct {
	Header       string
	Languages    []stats.Lang
	LanguagesLen int
}

func GenerateSVG(c *gin.Context) {
	header := "Languages"
	languages := []stats.Lang{
		{Name: "Java", Percent: 48.9},
		{Name: "JavaScript", Percent: 47.0},
		{Name: "CSS", Percent: 1.7},
		{Name: "Go", Percent: 1.6},
		{Name: "HTML", Percent: 0.8},
	}

	data := SVGData{
		Header:       header,
		Languages:    languages,
		LanguagesLen: len(languages),
	}

	svg := generateSVG(data)
	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svg)
}

func generateSVG(data SVGData) string {
	tmpl := template.New("template.svg").Funcs(template.FuncMap{
		"ge": func(a, b int) bool { return a >= b }})
	tmpl, err := tmpl.ParseFiles("internal/svg/template.svg")
	if err != nil {
		return "<svg><!-- template error --></svg>"
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "<svg><!-- execute error --></svg>"
	}

	return buf.String()
}
