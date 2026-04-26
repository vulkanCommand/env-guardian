package codebase

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type NamingMismatch struct {
	CodeKey string
	EnvKey  string
}

type ScanResult struct {
	UsedKeys          []string
	MissingInEnv      []string
	UnusedInEnv       []string
	NamingMismatches  []NamingMismatch
	ScannedFilesCount int
}

var envPatterns = []*regexp.Regexp{
	regexp.MustCompile(`os\.Getenv\("([A-Z0-9_]+)"\)`),
	regexp.MustCompile(`os\.LookupEnv\("([A-Z0-9_]+)"\)`),
	regexp.MustCompile(`process\.env\.([A-Z0-9_]+)`),
	regexp.MustCompile(`process\.env\["([A-Z0-9_]+)"\]`),
	regexp.MustCompile(`process\.env\['([A-Z0-9_]+)'\]`),
	regexp.MustCompile(`import\.meta\.env\.([A-Z0-9_]+)`),
	regexp.MustCompile(`os\.getenv\("([A-Z0-9_]+)"\)`),
	regexp.MustCompile(`os\.environ\["([A-Z0-9_]+)"\]`),
	regexp.MustCompile(`os\.environ\['([A-Z0-9_]+)'\]`),
}

var allowedExtensions = map[string]bool{
	".go":  true,
	".js":  true,
	".jsx": true,
	".ts":  true,
	".tsx": true,
	".py":  true,
	".mjs": true,
	".cjs": true,
}

func Run(rootDir string, envValues map[string]string) (ScanResult, error) {
	usedKeysSet := make(map[string]bool)
	scannedFilesCount := 0

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if shouldSkipDir(path, info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldSkipFile(path) {
			return nil
		}

		if !allowedExtensions[strings.ToLower(filepath.Ext(path))] {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		scannedFilesCount++

		for _, pattern := range envPatterns {
			matches := pattern.FindAllStringSubmatch(string(content), -1)
			for _, match := range matches {
				if len(match) > 1 {
					usedKeysSet[match[1]] = true
				}
			}
		}

		return nil
	})
	if err != nil {
		return ScanResult{}, err
	}

	usedKeys := make([]string, 0, len(usedKeysSet))
	for key := range usedKeysSet {
		usedKeys = append(usedKeys, key)
	}
	sort.Strings(usedKeys)

	envKeys := make([]string, 0, len(envValues))
	for key := range envValues {
		envKeys = append(envKeys, key)
	}
	sort.Strings(envKeys)

	normalizedEnvKeys := make(map[string][]string)
	for _, key := range envKeys {
		normalized := normalizeKey(key)
		normalizedEnvKeys[normalized] = append(normalizedEnvKeys[normalized], key)
	}

	matchedEnvKeys := make(map[string]bool)
	missingInEnv := []string{}
	namingMismatches := []NamingMismatch{}

	for _, codeKey := range usedKeys {
		if _, exists := envValues[codeKey]; exists {
			matchedEnvKeys[codeKey] = true
			continue
		}

		normalized := normalizeKey(codeKey)
		candidates := normalizedEnvKeys[normalized]

		if len(candidates) > 0 {
			sort.Strings(candidates)
			envKey := candidates[0]
			namingMismatches = append(namingMismatches, NamingMismatch{
				CodeKey: codeKey,
				EnvKey:  envKey,
			})
			matchedEnvKeys[envKey] = true
			continue
		}

		missingInEnv = append(missingInEnv, codeKey)
	}

	unusedInEnv := []string{}
	for _, envKey := range envKeys {
		if matchedEnvKeys[envKey] {
			continue
		}
		if !usedKeysSet[envKey] {
			unusedInEnv = append(unusedInEnv, envKey)
		}
	}

	sort.Strings(missingInEnv)
	sort.Strings(unusedInEnv)
	sort.Slice(namingMismatches, func(i, j int) bool {
		if namingMismatches[i].CodeKey == namingMismatches[j].CodeKey {
			return namingMismatches[i].EnvKey < namingMismatches[j].EnvKey
		}
		return namingMismatches[i].CodeKey < namingMismatches[j].CodeKey
	})

	return ScanResult{
		UsedKeys:          usedKeys,
		MissingInEnv:      missingInEnv,
		UnusedInEnv:       unusedInEnv,
		NamingMismatches:  namingMismatches,
		ScannedFilesCount: scannedFilesCount,
	}, nil
}

func shouldSkipDir(path string, name string) bool {
	switch name {
	case ".git", "node_modules", "dist", "build", "vendor", ".next", "coverage":
		return true
	}

	cleanPath := filepath.ToSlash(path)

	if strings.Contains(cleanPath, "/internal/codebase") {
		return true
	}

	return false
}

func shouldSkipFile(path string) bool {
	cleanPath := filepath.ToSlash(path)
	base := filepath.Base(path)

	if strings.HasSuffix(base, "_test.go") {
		return true
	}

	if cleanPath == "cmd/envguard/main.go" || strings.HasSuffix(cleanPath, "/cmd/envguard/main.go") {
		return true
	}

	return false
}

func normalizeKey(key string) string {
	var builder strings.Builder

	for _, ch := range key {
		if ch >= 'a' && ch <= 'z' {
			builder.WriteRune(ch - 32)
			continue
		}

		if (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			builder.WriteRune(ch)
		}
	}

	return builder.String()
}
