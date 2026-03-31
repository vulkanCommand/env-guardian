package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadTypeSchema(path string) (map[string]string, error) {
	schema := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return schema, nil
		}
		return nil, fmt.Errorf("failed to open type schema file: %w", err)
	}
	defer file.Close()

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
			return nil, fmt.Errorf("invalid schema format at line %d: %s", lineNumber, line)
		}

		key := strings.TrimSpace(parts[0])
		valueType := strings.TrimSpace(parts[1])

		if key == "" || valueType == "" {
			return nil, fmt.Errorf("invalid empty key/type at line %d", lineNumber)
		}

		schema[key] = valueType
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading schema file: %w", err)
	}

	return schema, nil
}
