package encryption

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	plaintext := []byte("PORT=3000\nJWT_SECRET=replace-me\n")
	key := "test-encryption-key"

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("expected no encryption error, got %v", err)
	}

	if !strings.HasPrefix(encrypted, formatPrefix+":") {
		t.Fatalf("expected encrypted value to use format prefix")
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("expected no decryption error, got %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Fatalf("expected decrypted content %q, got %q", string(plaintext), string(decrypted))
	}
}

func TestDecryptRejectsWrongKey(t *testing.T) {
	encrypted, err := Encrypt([]byte("PORT=3000\n"), "correct-key")
	if err != nil {
		t.Fatalf("expected no encryption error, got %v", err)
	}

	_, err = Decrypt(encrypted, "wrong-key")
	if err == nil {
		t.Fatal("expected wrong key to fail decryption")
	}
}

func TestDecryptRejectsInvalidFormat(t *testing.T) {
	_, err := Decrypt("not-encrypted", "key")
	if err == nil {
		t.Fatal("expected invalid format to fail")
	}
}

func TestEncryptFileDecryptFileRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, ".env")
	encryptedPath := filepath.Join(tempDir, ".env.enc")
	outputPath := filepath.Join(tempDir, ".env.dec")
	plaintext := "PORT=3000\nDEBUG=true\n"

	err := os.WriteFile(inputPath, []byte(plaintext), 0600)
	if err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	err = EncryptFile(inputPath, encryptedPath, "file-key")
	if err != nil {
		t.Fatalf("expected no encrypt file error, got %v", err)
	}

	err = DecryptFile(encryptedPath, outputPath, "file-key")
	if err != nil {
		t.Fatalf("expected no decrypt file error, got %v", err)
	}

	decrypted, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Fatalf("expected decrypted file %q, got %q", plaintext, string(decrypted))
	}
}

func TestEncryptFileRequiresKey(t *testing.T) {
	err := EncryptFile("missing.env", "out.env", "")
	if err == nil {
		t.Fatal("expected missing key error")
	}
}
