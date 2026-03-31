package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/vulkanCommand/env-guardian/internal/analyzer"
	"github.com/vulkanCommand/env-guardian/internal/linter"
	"github.com/vulkanCommand/env-guardian/internal/parser"
	"github.com/vulkanCommand/env-guardian/internal/validator"
	"github.com/vulkanCommand/env-guardian/internal/version"
)

func printHelp() {
	fmt.Println("envguard")
	fmt.Println("A CLI tool to validate, lint, and analyze environment variables.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  envguard help")
	fmt.Println("  envguard help validate")
	fmt.Println("  envguard help lint")
	fmt.Println("  envguard help analyze")
	fmt.Println("  envguard help doctor")
	fmt.Println("  envguard version")
	fmt.Println("  envguard validate")
	fmt.Println("  envguard validate --file .env.prod")
	fmt.Println("  envguard validate --file .env.prod --example .env.example.prod")
	fmt.Println("  envguard lint")
	fmt.Println("  envguard lint --file .env.prod")
	fmt.Println("  envguard analyze")
	fmt.Println("  envguard analyze --file .env.prod")
	fmt.Println("  envguard doctor")
	fmt.Println("  envguard doctor --file .env.prod --example .env.example.prod")
}

func printValidateHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard validate")
	fmt.Println("  envguard validate --file .env.prod")
	fmt.Println("  envguard validate --file .env.prod --example .env.example.prod")
	fmt.Println("")
	fmt.Println("Validation checks:")
	fmt.Println("  - missing keys compared to the example file")
	fmt.Println("  - duplicate keys in the target env file")
	fmt.Println("  - unused keys not present in the example file")
	fmt.Println("  - typed values from examples/.env.types")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Target env file to validate")
	fmt.Println("  --example   Example env file to compare against")
}

func printLintHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard lint")
	fmt.Println("  envguard lint --file .env.prod")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Target env file to lint")
}

func printAnalyzeHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard analyze")
	fmt.Println("  envguard analyze --file .env.prod")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Target env file to analyze")
}

func printDoctorHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard doctor")
	fmt.Println("  envguard doctor --file .env.prod --example .env.example.prod")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Target env file to inspect")
	fmt.Println("  --example   Example env file to compare against")
}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func handleHelpCommand(args []string) int {
	if len(args) == 0 {
		printHelp()
		return 0
	}

	if len(args) > 1 {
		fmt.Printf("Error: unexpected argument: %s\n", args[1])
		return 1
	}

	switch args[0] {
	case "validate":
		printValidateHelp()
		return 0
	case "lint":
		printLintHelp()
		return 0
	case "analyze":
		printAnalyzeHelp()
		return 0
	case "doctor":
		printDoctorHelp()
		return 0
	default:
		fmt.Printf("Error: unknown help topic: %s\n", args[0])
		return 1
	}
}

func getValidatePaths(args []string) (string, string, error) {
	envPath := ".env"
	examplePath := ".env.example"
	fileFlagSeen := false
	exampleFlagSeen := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--file" {
			if fileFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --file")
			}
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("missing value for --file")
			}

			next := args[i+1]
			if strings.HasPrefix(next, "--") {
				return "", "", fmt.Errorf("missing value for --file")
			}

			envPath = next
			fileFlagSeen = true
			i++
			continue
		}

		if strings.HasPrefix(arg, "--file=") {
			if fileFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --file")
			}

			value := strings.TrimPrefix(arg, "--file=")
			if value == "" {
				return "", "", fmt.Errorf("missing value for --file")
			}

			envPath = value
			fileFlagSeen = true
			continue
		}

		if arg == "--example" {
			if exampleFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --example")
			}
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("missing value for --example")
			}

			next := args[i+1]
			if strings.HasPrefix(next, "--") {
				return "", "", fmt.Errorf("missing value for --example")
			}

			examplePath = next
			exampleFlagSeen = true
			i++
			continue
		}

		if strings.HasPrefix(arg, "--example=") {
			if exampleFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --example")
			}

			value := strings.TrimPrefix(arg, "--example=")
			if value == "" {
				return "", "", fmt.Errorf("missing value for --example")
			}

			examplePath = value
			exampleFlagSeen = true
			continue
		}

		if strings.HasPrefix(arg, "--") {
			return "", "", fmt.Errorf("unknown flag: %s", arg)
		}

		return "", "", fmt.Errorf("unexpected argument: %s", arg)
	}

	return envPath, examplePath, nil
}

func getLintFilePath(args []string) (string, error) {
	filePath := ".env"
	fileFlagSeen := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--file" {
			if fileFlagSeen {
				return "", fmt.Errorf("duplicate flag: --file")
			}
			if i+1 >= len(args) {
				return "", fmt.Errorf("missing value for --file")
			}

			next := args[i+1]
			if strings.HasPrefix(next, "--") {
				return "", fmt.Errorf("missing value for --file")
			}

			filePath = next
			fileFlagSeen = true
			i++
			continue
		}

		if strings.HasPrefix(arg, "--file=") {
			if fileFlagSeen {
				return "", fmt.Errorf("duplicate flag: --file")
			}

			value := strings.TrimPrefix(arg, "--file=")
			if value == "" {
				return "", fmt.Errorf("missing value for --file")
			}

			filePath = value
			fileFlagSeen = true
			continue
		}

		if strings.HasPrefix(arg, "--") {
			return "", fmt.Errorf("unknown flag: %s", arg)
		}

		return "", fmt.Errorf("unexpected argument: %s", arg)
	}

	return filePath, nil
}

func getAnalyzeFilePath(args []string) (string, error) {
	filePath := ".env"
	fileFlagSeen := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--file" {
			if fileFlagSeen {
				return "", fmt.Errorf("duplicate flag: --file")
			}
			if i+1 >= len(args) {
				return "", fmt.Errorf("missing value for --file")
			}

			next := args[i+1]
			if strings.HasPrefix(next, "--") {
				return "", fmt.Errorf("missing value for --file")
			}

			filePath = next
			fileFlagSeen = true
			i++
			continue
		}

		if strings.HasPrefix(arg, "--file=") {
			if fileFlagSeen {
				return "", fmt.Errorf("duplicate flag: --file")
			}

			value := strings.TrimPrefix(arg, "--file=")
			if value == "" {
				return "", fmt.Errorf("missing value for --file")
			}

			filePath = value
			fileFlagSeen = true
			continue
		}

		if strings.HasPrefix(arg, "--") {
			return "", fmt.Errorf("unknown flag: %s", arg)
		}

		return "", fmt.Errorf("unexpected argument: %s", arg)
	}

	return filePath, nil
}

func getDoctorPaths(args []string) (string, string, error) {
	envPath := ".env"
	examplePath := ".env.example"
	fileFlagSeen := false
	exampleFlagSeen := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--file" {
			if fileFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --file")
			}
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("missing value for --file")
			}

			next := args[i+1]
			if strings.HasPrefix(next, "--") {
				return "", "", fmt.Errorf("missing value for --file")
			}

			envPath = next
			fileFlagSeen = true
			i++
			continue
		}

		if strings.HasPrefix(arg, "--file=") {
			if fileFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --file")
			}

			value := strings.TrimPrefix(arg, "--file=")
			if value == "" {
				return "", "", fmt.Errorf("missing value for --file")
			}

			envPath = value
			fileFlagSeen = true
			continue
		}

		if arg == "--example" {
			if exampleFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --example")
			}
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("missing value for --example")
			}

			next := args[i+1]
			if strings.HasPrefix(next, "--") {
				return "", "", fmt.Errorf("missing value for --example")
			}

			examplePath = next
			exampleFlagSeen = true
			i++
			continue
		}

		if strings.HasPrefix(arg, "--example=") {
			if exampleFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --example")
			}

			value := strings.TrimPrefix(arg, "--example=")
			if value == "" {
				return "", "", fmt.Errorf("missing value for --example")
			}

			examplePath = value
			exampleFlagSeen = true
			continue
		}

		if strings.HasPrefix(arg, "--") {
			return "", "", fmt.Errorf("unknown flag: %s", arg)
		}

		return "", "", fmt.Errorf("unexpected argument: %s", arg)
	}

	return envPath, examplePath, nil
}

func runValidate(envPath string, examplePath string) int {
	schema := map[string]string{}

	loadedSchema, err := parser.LoadTypeSchema("examples/.env.types")
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Error: could not read type schema file: %v\n", err)
			return 1
		}
	} else {
		schema = loadedSchema
	}

	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		fmt.Printf("Error: could not read %s\n", envPath)
		return 1
	}

	exampleFile, err := parser.ParseEnvFile(examplePath)
	if err != nil {
		fmt.Printf("Error: could not read %s\n", examplePath)
		return 1
	}

	result := validator.ValidateEnv(envFile, exampleFile, schema)

	sort.Strings(result.MissingKeys)
	sort.Strings(result.DuplicateKeys)
	sort.Strings(result.UnusedKeys)
	sort.Strings(result.InvalidTypeValues)

	fmt.Println("Env Validation Report")
	fmt.Println("---------------------")
	fmt.Printf("Target file: %s\n", envPath)
	fmt.Printf("Example file: %s\n\n", examplePath)

	errorCount := 0
	warningCount := 0

	for _, key := range result.MissingKeys {
		fmt.Printf("[ERROR] Missing key: %s\n", key)
		errorCount++
	}

	for _, key := range result.DuplicateKeys {
		fmt.Printf("[ERROR] Duplicate key: %s\n", key)
		errorCount++
	}

	for _, issue := range result.InvalidTypeValues {
		fmt.Printf("[ERROR] Invalid type: %s\n", issue)
		errorCount++
	}

	if errorCount == 0 && warningCount == 0 {
		fmt.Println("[PASS] Environment configuration looks good")
		fmt.Println("")
		fmt.Printf("Summary: %d error(s), %d warning(s)\n", errorCount, warningCount)
		return 0
	}

	if errorCount == 0 && warningCount > 0 {
		fmt.Println("[PASS] Environment configuration is valid with warnings")
		fmt.Println("")
		fmt.Printf("Summary: %d error(s), %d warning(s)\n", errorCount, warningCount)
		return 0
	}

	fmt.Println("")
	fmt.Printf("Summary: %d error(s), %d warning(s)\n", errorCount, warningCount)

	return 1
}

func runLint(envPath string) int {
	result, err := linter.Run(envPath)
	if err != nil {
		fmt.Printf("Error: could not read %s\n", envPath)
		return 1
	}

	fmt.Println("Env Lint Report")
	fmt.Println("---------------")
	fmt.Printf("Target file: %s\n\n", envPath)

	if len(result.InvalidLines) == 0 {
		fmt.Println("[PASS] Env syntax looks good")
		fmt.Println("")
		fmt.Printf("Summary: %d lint issue(s) found\n", len(result.InvalidLines))
		return 0
	}

	for _, issue := range result.InvalidLines {
		fmt.Printf("[ERROR] %s\n", issue)
	}

	fmt.Println("")
	fmt.Printf("Summary: %d lint issue(s) found\n", len(result.InvalidLines))

	return 1
}

func runAnalyze(envPath string) int {
	result, err := analyzer.Run(envPath)
	if err != nil {
		fmt.Printf("Error: could not read %s\n", envPath)
		return 1
	}

	sort.Strings(result.EmptyValues)
	sort.Strings(result.PotentialSecrets)

	fmt.Println("Env Analysis Report")
	fmt.Println("-------------------")
	fmt.Printf("Target file: %s\n", envPath)
	fmt.Printf("Total keys: %d\n", result.TotalKeys)

	if len(result.EmptyValues) > 0 {
		fmt.Println("\nEmpty values:")
		for _, key := range result.EmptyValues {
			fmt.Printf("- %s\n", key)
		}
	}

	if len(result.PotentialSecrets) > 0 {
		fmt.Println("\nPotential sensitive keys:")
		for _, key := range result.PotentialSecrets {
			fmt.Printf("- %s\n", key)
		}
	}

	fmt.Println("")
	fmt.Printf("Summary: %d empty value(s), %d potential sensitive key(s)\n", len(result.EmptyValues), len(result.PotentialSecrets))

	if len(result.EmptyValues) == 0 && len(result.PotentialSecrets) == 0 {
		fmt.Println("[PASS] Environment analysis found no issues")
	}

	return 0
}

func runDoctor(envPath string, examplePath string) int {
	envExists := true
	exampleExists := true
	missingInEnv := []string{}

	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		envExists = false
	}

	exampleFile, err := parser.ParseEnvFile(examplePath)
	if err != nil {
		exampleExists = false
	}

	if envExists && exampleExists {
		for key := range exampleFile.Values {
			if _, exists := envFile.Values[key]; !exists {
				missingInEnv = append(missingInEnv, key)
			}
		}
		sort.Strings(missingInEnv)
	}

	targetFileIssues := 0
	exampleFileIssues := 0

	if !envExists {
		targetFileIssues = 1
	}

	if !exampleExists {
		exampleFileIssues = 1
	}

	fmt.Println("Env Doctor Report")
	fmt.Println("------------------")
	fmt.Printf("Target file: %s\n", envPath)
	fmt.Printf("Example file: %s\n\n", examplePath)

	if envExists {
		fmt.Printf("[OK] %s file exists\n", envPath)
	} else {
		fmt.Printf("[ERROR] %s file missing\n", envPath)
	}

	if exampleExists {
		fmt.Printf("[OK] %s file exists\n", examplePath)
	} else {
		fmt.Printf("[WARNING] %s file missing\n", examplePath)
	}

	if len(missingInEnv) > 0 {
		fmt.Printf("\n[WARNING] Missing keys in %s (present in %s):\n", envPath, examplePath)
		for _, key := range missingInEnv {
			fmt.Printf("- %s\n", key)
		}
	}

	if len(missingInEnv) == 0 && envExists && exampleExists {
		fmt.Println("[PASS] Environment doctor checks passed")
		fmt.Println("")
		fmt.Printf("Summary: %d target file issue(s), %d example file issue(s), %d missing key(s)\n", targetFileIssues, exampleFileIssues, len(missingInEnv))
		return 0
	}

	fmt.Println("")
	fmt.Printf("Summary: %d target file issue(s), %d example file issue(s), %d missing key(s)\n", targetFileIssues, exampleFileIssues, len(missingInEnv))

	if len(missingInEnv) > 0 || !envExists {
		return 1
	}

	return 0
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	switch args[0] {
	case "version":
		fmt.Println(version.Version)
	case "validate":
		if hasHelpFlag(args[1:]) {
			printValidateHelp()
			os.Exit(0)
		}
		envPath, examplePath, err := getValidatePaths(args[1:])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runValidate(envPath, examplePath))
	case "lint":
		if hasHelpFlag(args[1:]) {
			printLintHelp()
			os.Exit(0)
		}
		envPath, err := getLintFilePath(args[1:])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runLint(envPath))
	case "analyze":
		if hasHelpFlag(args[1:]) {
			printAnalyzeHelp()
			os.Exit(0)
		}
		envPath, err := getAnalyzeFilePath(args[1:])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runAnalyze(envPath))
	case "doctor":
		if hasHelpFlag(args[1:]) {
			printDoctorHelp()
			os.Exit(0)
		}
		envPath, examplePath, err := getDoctorPaths(args[1:])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runDoctor(envPath, examplePath))
	case "help":
		os.Exit(handleHelpCommand(args[1:]))
	case "--help", "-h":
		printHelp()
	default:
		fmt.Printf("unknown command: %s\n\n", args[0])
		printHelp()
		os.Exit(1)
	}
}
