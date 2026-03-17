package analyzer

import (
	"os"
	"strings"
)

type AnalyzeResult struct {
	TotalKeys     int
	EmptyValues   []string
	PotentialSecrets []string
}

func Run(path string) (AnalyzeResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return AnalyzeResult{}, err
	}

	lines := strings.Split(string(data), "\n")

	result := AnalyzeResult{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		result.TotalKeys++

		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if value == "" {
			result.EmptyValues = append(result.EmptyValues, key)
		}

		if strings.Contains(strings.ToLower(key), "secret") ||
			strings.Contains(strings.ToLower(key), "password") ||
			strings.Contains(strings.ToLower(key), "token") {
			result.PotentialSecrets = append(result.PotentialSecrets, key)
		}
	}

	return result, nil
}