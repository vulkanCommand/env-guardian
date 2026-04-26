package runtimecheck

import (
	"bufio"
	"os"
	"regexp"
	"sort"
	"strings"
)

type DockerResult struct {
	ReferencedKeys []string
	MissingKeys    []string
}

var dockerReferencePattern = regexp.MustCompile(`\$\{?([A-Z][A-Z0-9_]*)\}?`)

func ValidateDockerfile(path string, envValues map[string]string) (DockerResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return DockerResult{}, err
	}
	defer file.Close()

	referencedSet := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		for _, key := range parseDockerInstructionKeys(line) {
			referencedSet[key] = true
		}

		for _, match := range dockerReferencePattern.FindAllStringSubmatch(line, -1) {
			if len(match) > 1 {
				referencedSet[match[1]] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return DockerResult{}, err
	}

	referencedKeys := make([]string, 0, len(referencedSet))
	for key := range referencedSet {
		referencedKeys = append(referencedKeys, key)
	}
	sort.Strings(referencedKeys)

	missingKeys := []string{}
	for _, key := range referencedKeys {
		if _, exists := envValues[key]; !exists {
			missingKeys = append(missingKeys, key)
		}
	}

	return DockerResult{
		ReferencedKeys: referencedKeys,
		MissingKeys:    missingKeys,
	}, nil
}

func parseDockerInstructionKeys(line string) []string {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return []string{}
	}

	instruction := strings.ToUpper(fields[0])

	switch instruction {
	case "ARG":
		return []string{parseDockerKey(fields[1])}
	case "ENV":
		keys := []string{}
		if len(fields) == 2 {
			keys = append(keys, parseDockerKey(fields[1]))
			return filterEmptyKeys(keys)
		}

		for _, field := range fields[1:] {
			if strings.Contains(field, "=") {
				keys = append(keys, parseDockerKey(field))
			}
		}

		if len(keys) == 0 {
			keys = append(keys, parseDockerKey(fields[1]))
		}

		return filterEmptyKeys(keys)
	default:
		return []string{}
	}
}

func parseDockerKey(value string) string {
	key := strings.SplitN(value, "=", 2)[0]
	key = strings.TrimSpace(key)

	if key == "" {
		return ""
	}

	for _, ch := range key {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return ""
		}
	}

	return key
}

func filterEmptyKeys(keys []string) []string {
	filtered := []string{}
	for _, key := range keys {
		if key != "" {
			filtered = append(filtered, key)
		}
	}

	return filtered
}
