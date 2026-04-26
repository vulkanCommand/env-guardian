package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

const formatPrefix = "ENVGUARD_V1"

func EncryptFile(inputPath string, outputPath string, key string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("missing ENVGUARD_KEY")
	}

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, []byte(encrypted+"\n"), 0600)
}

func DecryptFile(inputPath string, outputPath string, key string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("missing ENVGUARD_KEY")
	}

	encrypted, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	plaintext, err := Decrypt(strings.TrimSpace(string(encrypted)), key)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0600)
}

func Encrypt(plaintext []byte, key string) (string, error) {
	block, err := aes.NewCipher(deriveKey(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return strings.Join([]string{
		formatPrefix,
		base64.StdEncoding.EncodeToString(nonce),
		base64.StdEncoding.EncodeToString(ciphertext),
	}, ":"), nil
}

func Decrypt(value string, key string) ([]byte, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 3 || parts[0] != formatPrefix {
		return nil, fmt.Errorf("invalid encrypted file format")
	}

	nonce, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid encrypted nonce: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid encrypted payload: %w", err)
	}

	block, err := aes.NewCipher(deriveKey(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(nonce) != gcm.NonceSize() {
		return nil, fmt.Errorf("invalid encrypted nonce size")
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file")
	}

	return plaintext, nil
}

func deriveKey(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}
