package linter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type LintResult struct {
	InvalidLines []string
}

func Run(path string) (LintResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return LintResult{}, err
	}
	defer file.Close()

	result := LintResult{}
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if !strings.Contains(trimmed, "=") {
			result.InvalidLines = append(
				result.InvalidLines,
				fmt.Sprintf("line %d: missing '=' -> %s", lineNumber, line),
			)
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		key := strings.TrimSpace(parts[0])

		if key == "" {
			result.InvalidLines = append(
				result.InvalidLines,
				fmt.Sprintf("line %d: empty key -> %s", lineNumber, line),
			)
		}
	}

	if err := scanner.Err(); err != nil {
		return LintResult{}, err
	}

	return result, nil
}