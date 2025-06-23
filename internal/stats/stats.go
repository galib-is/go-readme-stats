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

type lang struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}

func FetchStats(username string, ignoredLangsPath string) []lang {
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

func calcStats(langTotals map[string]int) []lang {
	totalBytes := 0
	for _, bytes := range langTotals {
		totalBytes += bytes
	}

	var result []lang
	for langName, bytes := range langTotals {
		percent := float64(bytes) / float64(totalBytes) * 100
		percent = math.Round(percent*100) / 100

		result = append(result, lang{
			Name:    langName,
			Percent: percent,
		})
	}

	return result
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
