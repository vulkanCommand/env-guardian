package parser

import (
	"bufio"
	"os"
	"strings"

	"github.com/vulkanCommand/env-guardian/internal/models"
)

func ParseEnvFile(path string) (*models.EnvFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := &models.EnvFile{
		Values:     make(map[string]string),
		Duplicates: make(map[string]int),
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if _, exists := result.Values[key]; exists {
			result.Duplicates[key]++
		} else {
			result.Values[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}