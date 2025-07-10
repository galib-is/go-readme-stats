package stats

import (
	"testing"
)

func TestLoadLanguageColours(t *testing.T) {
	colours, err := loadLanguageColours("../../internal/data/colours.json")
	if err != nil {
		t.Errorf("loadLanguageColours() error = %v", err)
	}

	if len(colours) == 0 {
		t.Error("loadLanguageColours() returned empty map")
	}

	if goColour, exists := colours["Go"]; !exists {
		t.Error("Go language colour not found")
	} else if goColour != "#00ADD8" {
		t.Errorf("Go colour = %s, want #00ADD8", goColour)
	}
}

func TestLoadLanguageColours_InvalidFile(t *testing.T) {
	_, err := loadLanguageColours("nonexistent.json")
	if err == nil {
		t.Error("loadLanguageColours() expected error for nonexistent file")
	}
}

func TestReadIgnoredLanguages(t *testing.T) {
	ignored, err := readIgnoredLanguages("../../config/ignored_languages.json")
	if err != nil {
		t.Errorf("readIgnoredLanguages() error = %v", err)
	}

	if _, exists := ignored["HTML"]; !exists {
		t.Error("HTML should be in ignored languages")
	}

	if _, exists := ignored["CSS"]; !exists {
		t.Error("CSS should be in ignored languages")
	}
}

func TestReadIgnoredLanguages_InvalidFile(t *testing.T) {
	_, err := readIgnoredLanguages("nonexistent.json")
	if err == nil {
		t.Error("readIgnoredLanguages() expected error for nonexistent file")
	}
}

func TestAddLanguageColours(t *testing.T) {
	languages := []Lang{
		{Name: "Go", Percent: 50.0},
		{Name: "UnknownLanguage", Percent: 30.0},
		{Name: "Java", Percent: 20.0},
	}

	err := addLanguageColours(languages)
	if err != nil {
		t.Errorf("addLanguageColours() error = %v", err)
	}

	if languages[0].Colour == "" {
		t.Error("Go language should have a colour assigned")
	}

	if languages[1].Colour != defaultColour {
		t.Errorf("Unknown language colour = %s, want %s", languages[1].Colour, defaultColour)
	}
}
