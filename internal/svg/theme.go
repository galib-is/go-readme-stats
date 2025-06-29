package svg

const DefaultTheme = "dark"

type Theme struct {
	Background    string
	Text          string
	SecondaryText string
}

// themes contains predefined colour schemes for SVG generation.
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

// GetTheme returns the theme configuration for the given name.
// Falls back to DefaultTheme if the requested theme doesn't exist.
func GetTheme(name string) Theme {
	if theme, exists := themes[name]; exists {
		return theme
	}

	return themes[DefaultTheme]
}
