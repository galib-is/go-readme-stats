package stats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
)

type repo struct {
	Name string `json:"name"`
	Fork bool   `json:"fork"`
}

type Lang struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
	Colour  string
}

func FetchStats(username string, ignoredLangsPath string) []Lang {
	repos := fetchRepoNames(username)
	ignoredLangs := readIgnoredLanguages(ignoredLangsPath)

	langTotals := make(map[string]int)
	for _, repo := range repos {
		if repo.Fork {
			continue
		}

		langs := fetchRepoLanguages(username, repo.Name)
		for lang, bytes := range langs {
			if _, ignored := ignoredLangs[lang]; ignored {
				continue
			}

			langTotals[lang] += bytes
		}
	}
	return calcStats(langTotals)
}

func callAPI(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)

}

func fetchRepoNames(username string) []repo {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	body, err := callAPI(url)
	if err != nil {
		log.Fatalf("failed to fetch repos: %v", err)
	}

	var repos []repo
	if err = json.Unmarshal(body, &repos); err != nil {
		log.Fatalf("failed to unmarshal repos: %v", err)
	}

	return repos
}

func fetchRepoLanguages(username string, repoName string) map[string]int {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", username, repoName)
	body, err := callAPI(url)
	if err != nil {
		log.Fatalf("failed to fetch languages for %s: %v", repoName, err)
		return nil
	}

	var langs map[string]int
	if err := json.Unmarshal(body, &langs); err != nil {
		log.Printf("failed to unmarshal languages for %s: %v", repoName, err)
		return nil
	}

	return langs
}

func calcStats(langTotals map[string]int) []Lang {
	totalBytes := 0
	for _, bytes := range langTotals {
		totalBytes += bytes
	}

	colours := loadLanguageColors()

	var result []Lang
	for langName, bytes := range langTotals {
		percent := float64(bytes) / float64(totalBytes) * 100
		percent = math.Round(percent*10) / 10

		colour, exists := colours[langName]
		if !exists {
			colour = "#858585"
		}

		result = append(result, Lang{
			Name:    langName,
			Percent: percent,
			Colour:  colour,
		})
	}

	return result
}

func loadLanguageColors() map[string]string {
	data, err := os.ReadFile("internal/data/colours.json")
	if err != nil {
		log.Printf("Warning: Failed to read colors file: %v", err)
		return make(map[string]string)
	}

	var colors map[string]string
	if err := json.Unmarshal(data, &colors); err != nil {
		log.Printf("Warning: Failed to parse colors file: %v", err)
		return make(map[string]string)
	}

	return colors
}

func readIgnoredLanguages(filename string) map[string]struct{} {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", filename, err)
	}

	var langs []string
	if err := json.Unmarshal(data, &langs); err != nil {
		log.Fatalf("Failed to decode %s: %v", filename, err)
	}

	set := make(map[string]struct{})
	for _, lang := range langs {
		set[lang] = struct{}{}
	}

	return set
}
