package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTypeSchema_MissingFileReturnsEmptySchema(t *testing.T) {
	tempDir := t.TempDir()
	missingFilePath := filepath.Join(tempDir, ".env.types")

	schema, err := LoadTypeSchema(missingFilePath)
	if err != nil {
		t.Fatalf("expected no error for missing schema file, got %v", err)
	}

	if schema == nil {
		t.Fatal("expected empty schema map, got nil")
	}

	if len(schema) != 0 {
		t.Fatalf("expected empty schema map, got %d entries", len(schema))
	}
}

func TestLoadTypeSchema_ValidFile(t *testing.T) {
	tempDir := t.TempDir()
	schemaFilePath := filepath.Join(tempDir, ".env.types")

	content := "DEBUG=boolean\nPORT=number\nAPI_URL=url\n"
	err := os.WriteFile(schemaFilePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write schema file: %v", err)
	}

	schema, err := LoadTypeSchema(schemaFilePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(schema) != 3 {
		t.Fatalf("expected 3 schema entries, got %d", len(schema))
	}

	if schema["DEBUG"] != "boolean" {
		t.Fatalf("expected DEBUG=boolean, got %q", schema["DEBUG"])
	}

	if schema["PORT"] != "number" {
		t.Fatalf("expected PORT=number, got %q", schema["PORT"])
	}

	if schema["API_URL"] != "url" {
		t.Fatalf("expected API_URL=url, got %q", schema["API_URL"])
	}
}
