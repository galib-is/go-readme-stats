package svg

import (
	"bytes"
	"net/http"
	"text/template"

	"go-readme-stats/internal/stats"

	"github.com/gin-gonic/gin"
)

type SVGData struct {
	Height       float64
	Header       string
	Languages    []stats.Lang
	LanguagesLen int
}

func GenerateSVG(c *gin.Context) {
	header := "Languages"
	username := "galib-i"
	ignoredLangsPath := "ignored_languages.json"

	languages := stats.FetchStats(username, ignoredLangsPath)
	languageLen := len(languages)
	height := calculateSVGHeight(languageLen)

	data := SVGData{
		Height:       height,
		Header:       header,
		Languages:    languages,
		LanguagesLen: languageLen,
	}

	svg := generateSVG(data)
	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svg)
}

func calculateSVGHeight(languageCount int) float64 {
	switch {
	case languageCount <= 2:
		return 114.5
	case languageCount <= 4:
		return 134.5
	default:
		return 154.5
	}
}

func generateSVG(data SVGData) string {
	tmpl := template.New("template.svg").Funcs(template.FuncMap{
		"ge": func(a, b int) bool { return a >= b },
		"sumPrev": func(langs []stats.Lang, idx int) float64 {
			sum := 0.0
			for i := range idx {
				sum += langs[i].Percent
			}
			return sum
		},
	})
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
