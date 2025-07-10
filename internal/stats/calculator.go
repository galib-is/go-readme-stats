package stats

import (
	"fmt"
	"math"
	"sort"
)

const (
	maxVisibleLanguages = 6
	topLanguagesCount   = 5  // Number of individual languages before grouping into "Other"
	percentPrecision    = 10 // One decimal precision (e.g. 10.4%)
)

// calculateStats converts language byte counts to percentages.
// Languages are sorted by descending percentage, then ascending name.
// If more than maxVisibleLanguages exist, languages beyond topLanguagesCount are grouped into "Other".
func calculateStats(languageTotals map[string]int) []Lang {
	// Calculate total bytes across all languages
	totalBytes := 0
	for _, bytes := range languageTotals {
		totalBytes += bytes
	}

	var result []Lang
	for lang, bytes := range languageTotals {
		result = append(result, Lang{
			Name:    lang,
			Percent: float64(bytes) / float64(totalBytes) * 100,
		})
	}

	// Sort by percentage (desc), then by name (asc) for consistent ordering
	sort.Slice(result, func(i, j int) bool {
		if result[i].Percent == result[j].Percent {
			return result[i].Name < result[j].Name
		}

		return result[i].Percent > result[j].Percent
	})

	// Combine languages below topLanguagesCount into "Other"
	if len(result) > maxVisibleLanguages {
		var otherPercent float64
		for _, lang := range result[topLanguagesCount:] {
			otherPercent += lang.Percent
		}

		result = append(result[:topLanguagesCount], Lang{
			Name:    fmt.Sprintf("Other (%d)", len(result)-topLanguagesCount),
			Percent: otherPercent,
		})
	}

	for i := range result {
		result[i].Percent = roundPercent(result[i].Percent)
	}

	return result
}

func roundPercent(percent float64) float64 {
	return math.Round(percent*percentPrecision) / percentPrecision
}
