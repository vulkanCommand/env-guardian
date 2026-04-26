package security

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestScanEnvFileDetectsSecretValue(t *testing.T) {
	tempDir := t.TempDir()
	envPath := filepath.Join(tempDir, ".env")
	openAIKey := "sk-" + "abcdefghijklmnopqrstuvwxyz123456"

	err := os.WriteFile(envPath, []byte("OPENAI_API_KEY="+openAIKey+"\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	findings, err := ScanEnvFile(envPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(findings) == 0 {
		t.Fatal("expected secret finding")
	}
}

func TestScanEnvFileSkipsPlaceholderValues(t *testing.T) {
	tempDir := t.TempDir()
	envPath := filepath.Join(tempDir, ".env")

	err := os.WriteFile(envPath, []byte("JWT_SECRET=replace-me\nAPI_KEY=your_api_key\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	findings, err := ScanEnvFile(envPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(findings) != 0 {
		t.Fatalf("expected no findings for placeholders, got %v", findings)
	}
}

func TestScanRepositoryDetectsSecretPattern(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "app.js")
	stripeKey := "sk_live_" + "abcdefghijklmnopqrstuvwxyz"

	err := os.WriteFile(sourcePath, []byte("const key = \""+stripeKey+"\"\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	findings, err := ScanRepository(tempDir, filepath.Join(tempDir, ".env"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("expected 1 repository finding, got %d", len(findings))
	}

	if findings[0].Kind != "Stripe key" {
		t.Fatalf("expected Stripe key finding, got %q", findings[0].Kind)
	}
}

func TestIsGitTrackedDetectsTrackedEnvFile(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	tempDir := t.TempDir()
	runGit(t, tempDir, "init")

	envPath := filepath.Join(tempDir, ".env")
	err := os.WriteFile(envPath, []byte("PORT=3000\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	runGit(t, tempDir, "add", ".env")

	if !IsGitTracked(tempDir, envPath) {
		t.Fatal("expected .env to be tracked")
	}
}

func TestScanGitHistoryDetectsSecretPattern(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	tempDir := t.TempDir()
	runGit(t, tempDir, "init")
	runGit(t, tempDir, "config", "user.name", "test")
	runGit(t, tempDir, "config", "user.email", "test@example.com")

	secretPath := filepath.Join(tempDir, "secret.txt")
	awsKey := "AKIA" + "ABCDEFGHIJKLMNOP"
	err := os.WriteFile(secretPath, []byte(awsKey+"\n"), 0644)
	if err != nil {
		t.Fatalf("failed to write secret file: %v", err)
	}

	runGit(t, tempDir, "add", "secret.txt")
	runGit(t, tempDir, "commit", "-m", "add secret")

	findings, scanned := ScanGitHistory(tempDir)
	if !scanned {
		t.Fatal("expected git history scan to run")
	}

	if len(findings) == 0 {
		t.Fatal("expected git history finding")
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	command := exec.Command("git", append([]string{"-C", dir}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(output))
	}
}
