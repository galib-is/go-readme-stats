package stats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
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

func FetchStats(ignoredLangsPath string) ([]Lang, error) {
	repos, err := fetchRepoNames()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}

	ignoredLangs, err := readIgnoredLanguages(ignoredLangsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ignored languages: %w", err)
	}

	username, err := getUsername()
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}

	langTotals := make(map[string]int)
	for _, repo := range repos {
		if repo.Fork {
			continue
		}

		langs, err := fetchRepoLanguages(username, repo.Name)
		if err != nil {
			log.Printf("Warning: Failed to fetch languages for %s: %v", repo.Name, err)
			continue
		}

		for lang, bytes := range langs {
			if _, ignored := ignoredLangs[lang]; ignored {
				continue
			}

			langTotals[lang] += bytes
		}
	}

	stats := calculateStats(langTotals)
	if err := addLanguageColours(stats); err != nil {
		return nil, fmt.Errorf("failed to add colours: %w", err)
	}
	return stats, nil
}

func fetchRepoNames() ([]repository, error) {
	url := "https://api.github.com/user/repos"
	body, err := callAPI(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %w", err)
	}

	var repos []repository
	if err = json.Unmarshal(body, &repos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal repos: %w", err)
	}

	return repos, nil
}

func fetchRepoLanguages(username string, repoName string) (map[string]int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", username, repoName)
	body, err := callAPI(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch languages for %s: %w", repoName, err)
	}

	var languages map[string]int
	if err := json.Unmarshal(body, &languages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal languages for %s: %w", repoName, err)
	}

	return languages, nil
}

func callAPI(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to %s: %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status %d for URL %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)

}

func calculateStats(langTotals map[string]int) []Lang {
	totalBytes := 0
	for _, bytes := range langTotals {
		totalBytes += bytes
	}

	var result []Lang
	for langName, bytes := range langTotals {
		percent := float64(bytes) / float64(totalBytes) * 100
		percent = math.Round(percent*10) / 10

		result = append(result, Lang{
			Name:    langName,
			Percent: percent,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Percent == result[j].Percent {
			return result[i].Name < result[j].Name
		}

		return result[i].Percent > result[j].Percent
	})

	// Combine languages below top 5 into "Other"
	if len(result) > 6 {
		var otherPercent float64
		for _, lang := range result[5:] {
			otherPercent += lang.Percent
		}

		result = append(result[:5], Lang{
			Name:    fmt.Sprintf("Other (%d)", len(result)-5),
			Percent: math.Round(otherPercent*10) / 10,
		})
	}

	return result
}

func addLanguageColours(langs []Lang) error {
	colours, err := loadLanguageColours("internal/data/colours.json")
	if err != nil {
		log.Printf("Warning: Failed to load colours: %v", err)
		colours = make(map[string]string)
	}

	for i := range langs {
		colour, exists := colours[langs[i].Name]
		if !exists {
			colour = "#F0F6FC"
		}
		langs[i].Colour = colour
	}

	return nil
}

func loadLanguageColours(filepath string) (map[string]string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read colours file %s: %w", filepath, err)
	}

	var colours map[string]string
	if err := json.Unmarshal(data, &colours); err != nil {
		return nil, fmt.Errorf("failed to parse colours file %s: %w", filepath, err)
	}

	return colours, nil
}

func readIgnoredLanguages(filename string) (map[string]struct{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filename, err)
	}

	var langs []string
	if err := json.Unmarshal(data, &langs); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", filename, err)
	}

	set := make(map[string]struct{})
	for _, lang := range langs {
		set[lang] = struct{}{}
	}

	return set, nil
}

func getUsername() (string, error) {
	body, err := callAPI("https://api.github.com/user")
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}

	var user struct {
		Login string `json:"login"`
	}

	if err := json.Unmarshal(body, &user); err != nil {
		return "", fmt.Errorf("failed to parse user info: %w", err)
	}

	return user.Login, nil
}
