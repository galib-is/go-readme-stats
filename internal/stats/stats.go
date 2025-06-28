package stats

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	coloursFilePath = "internal/data/colours.json"
	defaultColour   = "#F0F6FC"
)

type repository struct {
	Name string `json:"name"`
	Fork bool   `json:"fork"`
}

type Lang struct {
	Name    string
	Percent float64
	Colour  string
}

func FetchStats(ignoredLanguagesPath string) ([]Lang, error) {
	repos, err := fetchRepoNames()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}

	ignoredLanguages, err := readIgnoredLanguages(ignoredLanguagesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ignored languages: %w", err)
	}

	username, err := getUsername()
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}

	languageTotals := make(map[string]int)
	for _, repo := range repos {
		if repo.Fork {
			continue
		}

		languages, err := fetchRepoLanguages(username, repo.Name)
		if err != nil {
			log.Printf("Warning: Failed to fetch languages for %s: %v", repo.Name, err)
			continue
		}

		for lang, bytes := range languages {
			if _, ignored := ignoredLanguages[lang]; ignored {
				continue
			}

			languageTotals[lang] += bytes
		}
	}

	stats := calculateStats(languageTotals)
	if err := addLanguageColours(stats); err != nil {
		return nil, fmt.Errorf("failed to add colours: %w", err)
	}
	return stats, nil
}

func addLanguageColours(languages []Lang) error {
	colours, err := loadLanguageColours(coloursFilePath)
	if err != nil {
		log.Printf("Warning: Failed to load colours: %v", err)
		colours = make(map[string]string)
	}

	for i := range languages {
		colour, exists := colours[languages[i].Name]
		if !exists {
			colour = defaultColour
		}
		languages[i].Colour = colour
	}

	return nil
}

func loadLanguageColours(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read colours file %s: %w", path, err)
	}

	var colours map[string]string
	if err := json.Unmarshal(data, &colours); err != nil {
		return nil, fmt.Errorf("failed to parse colours file %s: %w", path, err)
	}

	return colours, nil
}

func readIgnoredLanguages(path string) (map[string]struct{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var languages []string
	if err := json.Unmarshal(data, &languages); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", path, err)
	}

	set := make(map[string]struct{})
	for _, lang := range languages {
		set[lang] = struct{}{}
	}

	return set, nil
}
