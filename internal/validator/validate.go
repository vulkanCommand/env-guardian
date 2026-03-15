package validator

import (
	"github.com/vulkanCommand/env-guardian/internal/models"
)

type ValidationResult struct {
	MissingKeys   []string
	DuplicateKeys []string
}

func ValidateEnv(envFile *models.EnvFile, exampleFile *models.EnvFile) ValidationResult {
	result := ValidationResult{
		MissingKeys:   []string{},
		DuplicateKeys: []string{},
	}

	for key := range exampleFile.Values {
		if _, exists := envFile.Values[key]; !exists {
			result.MissingKeys = append(result.MissingKeys, key)
		}
	}

	for key := range envFile.Duplicates {
		result.DuplicateKeys = append(result.DuplicateKeys, key)
	}

	return result
}