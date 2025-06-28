package stats

import (
	"fmt"
	"math"
	"sort"
)

const (
	maxVisibleLanguages = 6
	topLanguagesCount   = 5
	percentPrecision    = 10
)

func calculateStats(languageTotals map[string]int) []Lang {
	totalBytes := 0
	for _, bytes := range languageTotals {
		totalBytes += bytes
	}

	var result []Lang
	for lang, bytes := range languageTotals {
		result = append(result, Lang{
			Name:    lang,
			Percent: roundPercent(float64(bytes) / float64(totalBytes) * 100),
		})
	}

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
			Percent: roundPercent(otherPercent),
		})
	}

	return result
}

func roundPercent(percent float64) float64 {
	return math.Round(percent*percentPrecision) / percentPrecision
}
