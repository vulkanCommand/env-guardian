package codebase

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRun_DetectsMissingAndUnusedEnvKeys(t *testing.T) {
	rootDir := t.TempDir()

	content := `package main

import "os"

func main() {
	_ = os.Getenv("APP_PORT")
	_ = os.LookupEnv("APP_DEBUG")
}
`

	err := os.WriteFile(filepath.Join(rootDir, "main.go"), []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	envValues := map[string]string{
		"APP_PORT":   "3000",
		"UNUSED_KEY": "value",
	}

	result, err := Run(rootDir, envValues)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedUsed := []string{"APP_DEBUG", "APP_PORT"}
	if !reflect.DeepEqual(result.UsedKeys, expectedUsed) {
		t.Fatalf("expected used keys %v, got %v", expectedUsed, result.UsedKeys)
	}

	expectedMissing := []string{"APP_DEBUG"}
	if !reflect.DeepEqual(result.MissingInEnv, expectedMissing) {
		t.Fatalf("expected missing keys %v, got %v", expectedMissing, result.MissingInEnv)
	}

	expectedUnused := []string{"UNUSED_KEY"}
	if !reflect.DeepEqual(result.UnusedInEnv, expectedUnused) {
		t.Fatalf("expected unused keys %v, got %v", expectedUnused, result.UnusedInEnv)
	}

	if len(result.NamingMismatches) != 0 {
		t.Fatalf("expected no naming mismatches, got %v", result.NamingMismatches)
	}

	if result.ScannedFilesCount != 1 {
		t.Fatalf("expected 1 scanned file, got %d", result.ScannedFilesCount)
	}
}

func TestRun_SupportsJavaScriptAndPythonPatterns(t *testing.T) {
	rootDir := t.TempDir()

	jsContent := `const apiUrl = process.env.API_URL
const apiKey = process.env["API_KEY"]
const mode = import.meta.env.MODE
`
	pyContent := `import os
port = os.getenv("PORT")
secret = os.environ["SECRET_KEY"]
`

	err := os.WriteFile(filepath.Join(rootDir, "app.js"), []byte(jsContent), 0644)
	if err != nil {
		t.Fatalf("failed to write js file: %v", err)
	}

	err = os.WriteFile(filepath.Join(rootDir, "app.py"), []byte(pyContent), 0644)
	if err != nil {
		t.Fatalf("failed to write python file: %v", err)
	}

	envValues := map[string]string{
		"API_URL":    "https://example.com",
		"API_KEY":    "abc",
		"MODE":       "development",
		"PORT":       "3000",
		"SECRET_KEY": "secret",
	}

	result, err := Run(rootDir, envValues)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedUsed := []string{"API_KEY", "API_URL", "MODE", "PORT", "SECRET_KEY"}
	if !reflect.DeepEqual(result.UsedKeys, expectedUsed) {
		t.Fatalf("expected used keys %v, got %v", expectedUsed, result.UsedKeys)
	}

	if len(result.MissingInEnv) != 0 {
		t.Fatalf("expected no missing keys, got %v", result.MissingInEnv)
	}

	if len(result.UnusedInEnv) != 0 {
		t.Fatalf("expected no unused keys, got %v", result.UnusedInEnv)
	}

	if len(result.NamingMismatches) != 0 {
		t.Fatalf("expected no naming mismatches, got %v", result.NamingMismatches)
	}

	if result.ScannedFilesCount != 2 {
		t.Fatalf("expected 2 scanned files, got %d", result.ScannedFilesCount)
	}
}

func TestRun_DetectsNamingMismatchAndDoesNotDoubleReport(t *testing.T) {
	rootDir := t.TempDir()

	content := `package main

import "os"

func main() {
	_ = os.Getenv("API_URL")
}
`

	err := os.WriteFile(filepath.Join(rootDir, "main.go"), []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	envValues := map[string]string{
		"APIURL": "https://example.com",
	}

	result, err := Run(rootDir, envValues)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedMismatches := []NamingMismatch{
		{CodeKey: "API_URL", EnvKey: "APIURL"},
	}
	if !reflect.DeepEqual(result.NamingMismatches, expectedMismatches) {
		t.Fatalf("expected naming mismatches %v, got %v", expectedMismatches, result.NamingMismatches)
	}

	if len(result.MissingInEnv) != 0 {
		t.Fatalf("expected no missing keys, got %v", result.MissingInEnv)
	}

	if len(result.UnusedInEnv) != 0 {
		t.Fatalf("expected no unused keys, got %v", result.UnusedInEnv)
	}
}

func TestRun_SkipsIgnoredDirectories(t *testing.T) {
	rootDir := t.TempDir()

	ignoredDir := filepath.Join(rootDir, "node_modules")
	err := os.MkdirAll(ignoredDir, 0755)
	if err != nil {
		t.Fatalf("failed to create ignored directory: %v", err)
	}

	err = os.WriteFile(filepath.Join(ignoredDir, "ignored.js"), []byte(`const token = process.env.SHOULD_NOT_APPEAR`), 0644)
	if err != nil {
		t.Fatalf("failed to write ignored file: %v", err)
	}

	result, err := Run(rootDir, map[string]string{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.UsedKeys) != 0 {
		t.Fatalf("expected no used keys, got %v", result.UsedKeys)
	}

	if result.ScannedFilesCount != 0 {
		t.Fatalf("expected 0 scanned files, got %d", result.ScannedFilesCount)
	}
}

func TestRun_SkipsTestFiles(t *testing.T) {
	rootDir := t.TempDir()

	err := os.WriteFile(filepath.Join(rootDir, "main_test.go"), []byte(`package main
import "os"
func TestSomething(t *testing.T) {
	_ = os.Getenv("SHOULD_NOT_APPEAR")
}`), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	result, err := Run(rootDir, map[string]string{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.UsedKeys) != 0 {
		t.Fatalf("expected no used keys, got %v", result.UsedKeys)
	}

	if result.ScannedFilesCount != 0 {
		t.Fatalf("expected 0 scanned files, got %d", result.ScannedFilesCount)
	}
}
