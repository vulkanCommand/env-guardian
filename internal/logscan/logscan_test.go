package logscan

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRunDetectsEnvAccessLoggedInSource(t *testing.T) {
	rootDir := t.TempDir()
	sourcePath := filepath.Join(rootDir, "main.go")

	content := `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Getenv("API_TOKEN"))
}
`

	err := os.WriteFile(sourcePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	result, err := Run(rootDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedKinds := []string{"env value logged"}
	actualKinds := []string{}
	for _, finding := range result.Findings {
		actualKinds = append(actualKinds, finding.Kind)
	}

	if !reflect.DeepEqual(actualKinds, expectedKinds) {
		t.Fatalf("expected findings %v, got %v", expectedKinds, actualKinds)
	}

	if result.ScannedFilesCount != 1 {
		t.Fatalf("expected 1 scanned file, got %d", result.ScannedFilesCount)
	}
}

func TestRunDetectsSecretsInLogFile(t *testing.T) {
	rootDir := t.TempDir()
	logPath := filepath.Join(rootDir, "app.log")
	stripeKey := "sk_live_" + "abcdefghijklmnopqrstuvwxyz"

	err := os.WriteFile(logPath, []byte("payment_key="+stripeKey+"\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	result, err := Run(rootDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result.Findings))
	}

	if result.Findings[0].Kind != "Stripe key" {
		t.Fatalf("expected Stripe key finding, got %q", result.Findings[0].Kind)
	}
}

func TestRunDetectsSensitiveLogValue(t *testing.T) {
	rootDir := t.TempDir()
	logPath := filepath.Join(rootDir, "app.log")

	err := os.WriteFile(logPath, []byte("JWT_SECRET=real-production-token\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	result, err := Run(rootDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result.Findings))
	}

	if result.Findings[0].Kind != "sensitive value in log" {
		t.Fatalf("expected sensitive value finding, got %q", result.Findings[0].Kind)
	}
}

func TestRunSkipsIgnoredDirectories(t *testing.T) {
	rootDir := t.TempDir()
	ignoredDir := filepath.Join(rootDir, "node_modules")

	err := os.MkdirAll(ignoredDir, 0755)
	if err != nil {
		t.Fatalf("failed to create ignored directory: %v", err)
	}

	err = os.WriteFile(filepath.Join(ignoredDir, "app.log"), []byte("JWT_SECRET=real-production-token\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write ignored log file: %v", err)
	}

	result, err := Run(rootDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Findings) != 0 {
		t.Fatalf("expected no findings, got %v", result.Findings)
	}

	if result.ScannedFilesCount != 0 {
		t.Fatalf("expected 0 scanned files, got %d", result.ScannedFilesCount)
	}
}
