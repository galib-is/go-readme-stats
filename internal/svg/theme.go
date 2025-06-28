package svg

const DefaultTheme = "dark"

type Theme struct {
	Background    string
	Text          string
	SecondaryText string
}

var themes = map[string]Theme{
	"dark": {
		Background:    "#0D1117",
		Text:          "#F0F6FC",
		SecondaryText: "#9198A1",
	},
	"light": {
		Background:    "#FFFFFF",
		Text:          "#1F2328",
		SecondaryText: "#59636E",
	},
}

func GetTheme(name string) Theme {
	if theme, exists := themes[name]; exists {
		return theme
	}

	return themes[DefaultTheme]
}
