package validator

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/vulkanCommand/env-guardian/internal/models"
)

type ValidationResult struct {
	MissingKeys       []string
	DuplicateKeys     []string
	UnusedKeys        []string
	InvalidTypeValues []string
}

func ValidateEnv(envFile *models.EnvFile, exampleFile *models.EnvFile, schema map[string]string) ValidationResult {
	result := ValidationResult{
		MissingKeys:       []string{},
		DuplicateKeys:     []string{},
		UnusedKeys:        []string{},
		InvalidTypeValues: []string{},
	}

	for key := range exampleFile.Values {
		if _, exists := envFile.Values[key]; !exists {
			result.MissingKeys = append(result.MissingKeys, key)
		}
	}

	for key := range envFile.Duplicates {
		result.DuplicateKeys = append(result.DuplicateKeys, key)
	}

	for key := range envFile.Values {
		if _, exists := exampleFile.Values[key]; !exists {
			result.UnusedKeys = append(result.UnusedKeys, key)
		}
	}

	for key, expectedType := range schema {
		value, exists := envFile.Values[key]
		if !exists {
			continue
		}

		switch strings.ToLower(strings.TrimSpace(expectedType)) {
		case "boolean":
			if !isBoolean(value) {
				result.InvalidTypeValues = append(
					result.InvalidTypeValues,
					fmt.Sprintf("%s expected boolean but got %q", key, value),
				)
			}
		case "number":
			if !isNumber(value) {
				result.InvalidTypeValues = append(
					result.InvalidTypeValues,
					fmt.Sprintf("%s expected number but got %q", key, value),
				)
			}
		case "url":
			if !isURL(value) {
				result.InvalidTypeValues = append(
					result.InvalidTypeValues,
					fmt.Sprintf("%s expected url but got %q", key, value),
				)
			}
		}
	}

	return result
}

func isBoolean(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	return normalized == "true" || normalized == "false"
}

func isNumber(value string) bool {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return false
	}

	_, err := strconv.ParseFloat(normalized, 64)
	return err == nil
}

func isURL(value string) bool {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return false
	}

	parsed, err := url.ParseRequestURI(normalized)
	if err != nil {
		return false
	}

	return parsed.Scheme != "" && parsed.Host != ""
}
