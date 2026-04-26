package logscan

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/vulkanCommand/env-guardian/internal/security"
)

type Finding struct {
	Location string
	Kind     string
}

type Result struct {
	Findings          []Finding
	ScannedFilesCount int
}

var sourceExtensions = map[string]bool{
	".go":  true,
	".js":  true,
	".jsx": true,
	".ts":  true,
	".tsx": true,
	".py":  true,
}

func Run(rootDir string) (Result, error) {
	findings := []Finding{}
	scannedFilesCount := 0

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if !shouldScanFile(path, info) {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		scannedFilesCount++

		if isLogFile(path) {
			findings = append(findings, scanLogFile(path, content)...)
		} else {
			findings = append(findings, scanSourceFile(path, content)...)
		}

		return nil
	})
	if err != nil {
		return Result{}, err
	}

	return Result{
		Findings:          sortFindings(dedupeFindings(findings)),
		ScannedFilesCount: scannedFilesCount,
	}, nil
}

func scanSourceFile(path string, content []byte) []Finding {
	findings := []Finding{}
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		location := path + ":" + strconv.Itoa(lineNumber)

		if containsLoggingCall(line) && containsEnvAccess(line) {
			findings = append(findings, Finding{Location: location, Kind: "env value logged"})
		}
	}

	return findings
}

func scanLogFile(path string, content []byte) []Finding {
	findings := []Finding{}
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		location := path + ":" + strconv.Itoa(lineNumber)

		for _, kind := range security.DetectSecretKinds(line) {
			findings = append(findings, Finding{Location: location, Kind: kind})
		}

		if looksLikeSensitiveLogValue(line) {
			findings = append(findings, Finding{Location: location, Kind: "sensitive value in log"})
		}
	}

	return findings
}

func shouldScanFile(path string, info os.FileInfo) bool {
	if info.Size() > 2*1024*1024 {
		return false
	}

	base := filepath.Base(path)
	if base == "envguard" || strings.HasSuffix(base, ".exe") {
		return false
	}

	if strings.HasSuffix(base, "_test.go") {
		return false
	}

	if isLogFile(path) {
		return true
	}

	return sourceExtensions[strings.ToLower(filepath.Ext(path))]
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git", "node_modules", "dist", "build", "vendor", ".next", "coverage", "bin":
		return true
	default:
		return false
	}
}

func isLogFile(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".log"
}

func containsLoggingCall(line string) bool {
	lower := strings.ToLower(line)
	patterns := []string{
		"fmt.print",
		"log.print",
		"console.log",
		"console.error",
		"console.warn",
		"print(",
		"logging.",
		"logger.",
	}

	for _, pattern := range patterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	return false
}

func containsEnvAccess(line string) bool {
	codeOnly := stripQuotedText(line)
	patterns := []string{
		"os.Getenv(",
		"os.LookupEnv(",
		"process.env",
		"import.meta.env",
		"os.getenv(",
		"os.environ",
	}

	for _, pattern := range patterns {
		if strings.Contains(codeOnly, pattern) {
			return true
		}
	}

	return false
}

func stripQuotedText(line string) string {
	var builder strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false

	for _, ch := range line {
		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}

		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if inSingleQuote || inDoubleQuote {
			continue
		}

		builder.WriteRune(ch)
	}

	return builder.String()
}

func containsSensitiveKeyword(line string) bool {
	normalized := strings.ToUpper(line)
	keywords := []string{
		"SECRET",
		"PASSWORD",
		"TOKEN",
		"API_KEY",
		"PRIVATE_KEY",
		"CREDENTIAL",
		"ACCESS_KEY",
	}

	for _, keyword := range keywords {
		if strings.Contains(normalized, keyword) {
			return true
		}
	}

	return false
}

func looksLikeSensitiveLogValue(line string) bool {
	if !containsSensitiveKeyword(line) {
		return false
	}

	parts := strings.FieldsFunc(line, func(ch rune) bool {
		return ch == '=' || ch == ':'
	})
	if len(parts) < 2 {
		return false
	}

	value := strings.TrimSpace(parts[len(parts)-1])
	value = strings.Trim(value, `"'`)

	if len(value) < 8 {
		return false
	}

	return !isPlaceholderValue(value)
}

func isPlaceholderValue(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	placeholders := []string{
		"replace-me",
		"change-me",
		"placeholder",
		"example",
		"dummy",
		"test",
		"your_",
		"your-",
	}

	for _, placeholder := range placeholders {
		if strings.Contains(normalized, placeholder) {
			return true
		}
	}

	return normalized == ""
}

func dedupeFindings(findings []Finding) []Finding {
	seen := make(map[string]bool)
	deduped := []Finding{}

	for _, finding := range findings {
		key := finding.Location + "|" + finding.Kind
		if seen[key] {
			continue
		}
		seen[key] = true
		deduped = append(deduped, finding)
	}

	return deduped
}

func sortFindings(findings []Finding) []Finding {
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Location == findings[j].Location {
			return findings[i].Kind < findings[j].Kind
		}
		return findings[i].Location < findings[j].Location
	})

	return findings
}
