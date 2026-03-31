package validator

import (
	"testing"

	"github.com/vulkanCommand/env-guardian/internal/models"
)

func TestValidateEnv_InvalidBooleanType(t *testing.T) {
	envFile := &models.EnvFile{
		Values: map[string]string{
			"DEBUG": "yes",
		},
		Duplicates: map[string]int{},
	}

	exampleFile := &models.EnvFile{
		Values: map[string]string{
			"DEBUG": "",
		},
		Duplicates: map[string]int{},
	}

	schema := map[string]string{
		"DEBUG": "boolean",
	}

	result := ValidateEnv(envFile, exampleFile, schema)

	if len(result.InvalidTypeValues) != 1 {
		t.Fatalf("expected 1 invalid type error, got %d", len(result.InvalidTypeValues))
	}

	expected := `DEBUG expected boolean but got "yes"`
	if result.InvalidTypeValues[0] != expected {
		t.Fatalf("expected %q, got %q", expected, result.InvalidTypeValues[0])
	}
}

func TestValidateEnv_InvalidNumberType(t *testing.T) {
	envFile := &models.EnvFile{
		Values: map[string]string{
			"PORT": "yes",
		},
		Duplicates: map[string]int{},
	}

	exampleFile := &models.EnvFile{
		Values: map[string]string{
			"PORT": "",
		},
		Duplicates: map[string]int{},
	}

	schema := map[string]string{
		"PORT": "number",
	}

	result := ValidateEnv(envFile, exampleFile, schema)

	if len(result.InvalidTypeValues) != 1 {
		t.Fatalf("expected 1 invalid type error, got %d", len(result.InvalidTypeValues))
	}

	expected := `PORT expected number but got "yes"`
	if result.InvalidTypeValues[0] != expected {
		t.Fatalf("expected %q, got %q", expected, result.InvalidTypeValues[0])
	}
}

func TestValidateEnv_InvalidURLType(t *testing.T) {
	envFile := &models.EnvFile{
		Values: map[string]string{
			"API_URL": "not-a-url",
		},
		Duplicates: map[string]int{},
	}

	exampleFile := &models.EnvFile{
		Values: map[string]string{
			"API_URL": "",
		},
		Duplicates: map[string]int{},
	}

	schema := map[string]string{
		"API_URL": "url",
	}

	result := ValidateEnv(envFile, exampleFile, schema)

	if len(result.InvalidTypeValues) != 1 {
		t.Fatalf("expected 1 invalid type error, got %d", len(result.InvalidTypeValues))
	}

	expected := `API_URL expected url but got "not-a-url"`
	if result.InvalidTypeValues[0] != expected {
		t.Fatalf("expected %q, got %q", expected, result.InvalidTypeValues[0])
	}
}

func TestValidateEnv_ValidTypedValues(t *testing.T) {
	envFile := &models.EnvFile{
		Values: map[string]string{
			"DEBUG":   "true",
			"PORT":    "3000",
			"API_URL": "https://example.com",
		},
		Duplicates: map[string]int{},
	}

	exampleFile := &models.EnvFile{
		Values: map[string]string{
			"DEBUG":   "",
			"PORT":    "",
			"API_URL": "",
		},
		Duplicates: map[string]int{},
	}

	schema := map[string]string{
		"DEBUG":   "boolean",
		"PORT":    "number",
		"API_URL": "url",
	}

	result := ValidateEnv(envFile, exampleFile, schema)

	if len(result.InvalidTypeValues) != 0 {
		t.Fatalf("expected 0 invalid type errors, got %d", len(result.InvalidTypeValues))
	}
}

func TestValidateEnv_EmptySchemaSkipsTypeValidation(t *testing.T) {
	envFile := &models.EnvFile{
		Values: map[string]string{
			"DEBUG": "not-a-boolean",
			"PORT":  "not-a-number",
		},
		Duplicates: map[string]int{},
	}

	exampleFile := &models.EnvFile{
		Values: map[string]string{
			"DEBUG": "",
			"PORT":  "",
		},
		Duplicates: map[string]int{},
	}

	schema := map[string]string{}

	result := ValidateEnv(envFile, exampleFile, schema)

	if len(result.InvalidTypeValues) != 0 {
		t.Fatalf("expected 0 invalid type errors when schema is empty, got %d", len(result.InvalidTypeValues))
	}
}
