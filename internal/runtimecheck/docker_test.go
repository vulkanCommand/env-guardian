package runtimecheck

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestValidateDockerfileDetectsMissingKeys(t *testing.T) {
	tempDir := t.TempDir()
	dockerfilePath := filepath.Join(tempDir, "Dockerfile")

	content := `FROM alpine
ARG API_URL
ENV PORT=3000 DEBUG=true
RUN echo ${DATABASE_URL}
`

	err := os.WriteFile(dockerfilePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write Dockerfile: %v", err)
	}

	envValues := map[string]string{
		"API_URL": "https://example.com",
		"PORT":    "3000",
	}

	result, err := ValidateDockerfile(dockerfilePath, envValues)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedReferenced := []string{"API_URL", "DATABASE_URL", "DEBUG", "PORT"}
	if !reflect.DeepEqual(result.ReferencedKeys, expectedReferenced) {
		t.Fatalf("expected referenced keys %v, got %v", expectedReferenced, result.ReferencedKeys)
	}

	expectedMissing := []string{"DATABASE_URL", "DEBUG"}
	if !reflect.DeepEqual(result.MissingKeys, expectedMissing) {
		t.Fatalf("expected missing keys %v, got %v", expectedMissing, result.MissingKeys)
	}
}

func TestValidateDockerfileHandlesENVKeyValueSyntax(t *testing.T) {
	tempDir := t.TempDir()
	dockerfilePath := filepath.Join(tempDir, "Dockerfile")

	content := `FROM alpine
ENV APP_ENV production
`

	err := os.WriteFile(dockerfilePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write Dockerfile: %v", err)
	}

	result, err := ValidateDockerfile(dockerfilePath, map[string]string{
		"APP_ENV": "production",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedReferenced := []string{"APP_ENV"}
	if !reflect.DeepEqual(result.ReferencedKeys, expectedReferenced) {
		t.Fatalf("expected referenced keys %v, got %v", expectedReferenced, result.ReferencedKeys)
	}

	if len(result.MissingKeys) != 0 {
		t.Fatalf("expected no missing keys, got %v", result.MissingKeys)
	}
}
