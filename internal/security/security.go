package security

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Finding struct {
	Location string
	Kind     string
}

type Result struct {
	EnvFindings        []Finding
	RepositoryFindings []Finding
	HistoryFindings    []Finding
	EnvFileTracked     bool
	HistoryScanned     bool
}

type secretPattern struct {
	Kind    string
	Pattern *regexp.Regexp
}

var secretPatterns = []secretPattern{
	{Kind: "AWS access key", Pattern: regexp.MustCompile(`\b(?:AKIA|ASIA)[A-Z0-9]{16}\b`)},
	{Kind: "OpenAI API key", Pattern: regexp.MustCompile(`\bsk-(?:proj-)?[A-Za-z0-9_-]{20,}\b`)},
	{Kind: "Stripe key", Pattern: regexp.MustCompile(`\b(?:sk|rk)_(?:live|test)_[A-Za-z0-9]{16,}\b`)},
	{Kind: "GitHub token", Pattern: regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9_]{20,}\b`)},
	{Kind: "Slack token", Pattern: regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9-]{20,}\b`)},
	{Kind: "private key", Pattern: regexp.MustCompile(`-----BEGIN (?:RSA |EC |OPENSSH |DSA )?PRIVATE KEY-----`)},
}

var securityFileExtensions = map[string]bool{
	".env":    true,
	".go":     true,
	".js":     true,
	".jsx":    true,
	".ts":     true,
	".tsx":    true,
	".py":     true,
	".json":   true,
	".yaml":   true,
	".yml":    true,
	".toml":   true,
	".md":     true,
	".txt":    true,
	".sh":     true,
	".ps1":    true,
	".docker": true,
}

func DetectSecretKinds(value string) []string {
	kinds := []string{}

	for _, pattern := range secretPatterns {
		if pattern.Pattern.MatchString(value) {
			kinds = append(kinds, pattern.Kind)
		}
	}

	sort.Strings(kinds)

	return kinds
}

func Run(rootDir string, envPath string) (Result, error) {
	result := Result{}

	result.EnvFileTracked = IsGitTracked(rootDir, envPath)

	envFindings, err := ScanEnvFile(envPath)
	if err != nil {
		return Result{}, err
	}
	result.EnvFindings = envFindings

	repositoryFindings, err := ScanRepository(rootDir, envPath)
	if err != nil {
		return Result{}, err
	}
	result.RepositoryFindings = repositoryFindings

	historyFindings, historyScanned := ScanGitHistory(rootDir)
	result.HistoryFindings = historyFindings
	result.HistoryScanned = historyScanned

	return result, nil
}

func ScanEnvFile(path string) ([]Finding, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	findings := []Finding{}
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := normalizeValue(parts[1])

		if isPlaceholderValue(value) {
			continue
		}

		location := path + ":" + key

		patternMatched := false
		for _, kind := range DetectSecretKinds(value) {
			findings = append(findings, Finding{Location: location, Kind: kind})
			patternMatched = true
		}

		if !patternMatched && isSensitiveKey(key) && looksSensitiveValue(value) {
			findings = append(findings, Finding{Location: location, Kind: "sensitive env value"})
		}

		if hasURLPassword(value) {
			findings = append(findings, Finding{Location: location, Kind: "credential in URL"})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sortFindings(dedupeFindings(findings)), nil
}

func ScanRepository(rootDir string, envPath string) ([]Finding, error) {
	findings := []Finding{}
	rootDir = filepath.Clean(rootDir)
	envAbsPath, _ := filepath.Abs(envPath)

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

		if shouldSkipSecurityFile(path, info, envAbsPath) {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		findings = append(findings, scanContent(path, content)...)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return sortFindings(dedupeFindings(findings)), nil
}

func ScanGitHistory(rootDir string) ([]Finding, bool) {
	if !isGitRepository(rootDir) {
		return []Finding{}, false
	}

	command := exec.Command("git", "-C", rootDir, "log", "-p", "--all", "--no-ext-diff")
	output, err := command.Output()
	if err != nil {
		return []Finding{}, false
	}

	findings := []Finding{}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	currentCommit := ""

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "commit ") {
			currentCommit = strings.TrimSpace(strings.TrimPrefix(line, "commit "))
			if len(currentCommit) > 12 {
				currentCommit = currentCommit[:12]
			}
			continue
		}

		for _, kind := range DetectSecretKinds(line) {
			location := "git history"
			if currentCommit != "" {
				location += ":" + currentCommit
			}
			findings = append(findings, Finding{Location: location, Kind: kind})
		}
	}

	return sortFindings(dedupeFindings(findings)), true
}

func IsGitTracked(rootDir string, path string) bool {
	if !isGitRepository(rootDir) {
		return false
	}

	relativePath := gitRelativePath(rootDir, path)
	command := exec.Command("git", "-C", rootDir, "ls-files", "--error-unmatch", "--", relativePath)
	return command.Run() == nil
}

func scanContent(path string, content []byte) []Finding {
	findings := []Finding{}
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		for _, kind := range DetectSecretKinds(line) {
			findings = append(findings, Finding{
				Location: path + ":" + strconv.Itoa(lineNumber),
				Kind:     kind,
			})
		}
	}

	return findings
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git", "node_modules", "dist", "build", "vendor", ".next", "coverage", "bin":
		return true
	default:
		return false
	}
}

func shouldSkipSecurityFile(path string, info os.FileInfo, envAbsPath string) bool {
	absPath, err := filepath.Abs(path)
	if err == nil && filepath.Clean(absPath) == filepath.Clean(envAbsPath) {
		return true
	}

	if info.Size() > 2*1024*1024 {
		return true
	}

	base := filepath.Base(path)
	if base == "envguard" || strings.HasSuffix(base, ".exe") {
		return true
	}

	ext := strings.ToLower(filepath.Ext(path))
	if securityFileExtensions[ext] {
		return false
	}

	if strings.HasPrefix(base, ".env") {
		return false
	}

	if base == "Dockerfile" {
		return false
	}

	return true
}

func isGitRepository(rootDir string) bool {
	command := exec.Command("git", "-C", rootDir, "rev-parse", "--is-inside-work-tree")
	return command.Run() == nil
}

func gitRelativePath(rootDir string, path string) string {
	if filepath.IsAbs(path) {
		rootAbsPath, err := filepath.Abs(rootDir)
		if err != nil {
			return filepath.ToSlash(filepath.Clean(path))
		}

		relativePath, err := filepath.Rel(rootAbsPath, path)
		if err == nil {
			return filepath.ToSlash(relativePath)
		}
	}

	return filepath.ToSlash(filepath.Clean(path))
}

func isSensitiveKey(key string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(key))
	sensitiveParts := []string{
		"SECRET",
		"PASSWORD",
		"TOKEN",
		"API_KEY",
		"PRIVATE_KEY",
		"CREDENTIAL",
		"ACCESS_KEY",
	}

	for _, part := range sensitiveParts {
		if strings.Contains(normalized, part) {
			return true
		}
	}

	return false
}

func looksSensitiveValue(value string) bool {
	if len(value) < 8 {
		return false
	}

	if isPlaceholderValue(value) {
		return false
	}

	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "true" || normalized == "false" {
		return false
	}

	return true
}

func hasURLPassword(value string) bool {
	normalized := strings.TrimSpace(value)
	if !strings.Contains(normalized, "://") || !strings.Contains(normalized, "@") {
		return false
	}

	prefix := strings.SplitN(normalized, "@", 2)[0]
	return strings.Contains(prefix, ":")
}

func normalizeValue(value string) string {
	trimmed := strings.TrimSpace(value)
	trimmed = strings.Trim(trimmed, `"`)
	trimmed = strings.Trim(trimmed, `'`)
	return trimmed
}

func isPlaceholderValue(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return true
	}

	placeholders := []string{
		"change-me",
		"changeme",
		"replace-me",
		"replace_me",
		"placeholder",
		"example",
		"dummy",
		"test",
		"todo",
		"your_",
		"your-",
		"<",
		">",
	}

	for _, placeholder := range placeholders {
		if strings.Contains(normalized, placeholder) {
			return true
		}
	}

	return false
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
