package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/vulkanCommand/env-guardian/internal/analyzer"
	"github.com/vulkanCommand/env-guardian/internal/codebase"
	"github.com/vulkanCommand/env-guardian/internal/encryption"
	"github.com/vulkanCommand/env-guardian/internal/linter"
	"github.com/vulkanCommand/env-guardian/internal/logscan"
	"github.com/vulkanCommand/env-guardian/internal/models"
	"github.com/vulkanCommand/env-guardian/internal/parser"
	"github.com/vulkanCommand/env-guardian/internal/security"
	"github.com/vulkanCommand/env-guardian/internal/validator"
	"github.com/vulkanCommand/env-guardian/internal/version"
)

func printHelp() {
	fmt.Println("envguard")
	fmt.Println("A CLI tool to validate, lint, analyze, and secure environment variables.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  envguard help")
	fmt.Println("  envguard help validate")
	fmt.Println("  envguard help lint")
	fmt.Println("  envguard help analyze")
	fmt.Println("  envguard help doctor")
	fmt.Println("  envguard help scan-code")
	fmt.Println("  envguard help security")
	fmt.Println("  envguard help log-scan")
	fmt.Println("  envguard help encrypt")
	fmt.Println("  envguard help decrypt")
	fmt.Println("  envguard version")
	fmt.Println("  envguard validate")
	fmt.Println("  envguard validate --all")
	fmt.Println("  envguard validate --file .env.prod")
	fmt.Println("  envguard validate --file .env.prod --example .env.example.prod")
	fmt.Println("  envguard lint")
	fmt.Println("  envguard lint --file .env.prod")
	fmt.Println("  envguard analyze")
	fmt.Println("  envguard analyze --file .env.prod")
	fmt.Println("  envguard doctor")
	fmt.Println("  envguard doctor --file .env.prod --example .env.example.prod")
	fmt.Println("  envguard scan-code")
	fmt.Println("  envguard scan-code --dir .")
	fmt.Println("  envguard scan-code --dir . --file .env.prod")
	fmt.Println("  envguard security")
	fmt.Println("  envguard security --dir . --file .env.prod")
	fmt.Println("  envguard log-scan")
	fmt.Println("  envguard log-scan --dir .")
	fmt.Println("  envguard encrypt")
	fmt.Println("  envguard encrypt --file .env.prod --out .env.prod.enc")
	fmt.Println("  envguard decrypt")
	fmt.Println("  envguard decrypt --file .env.prod.enc --out .env.prod")
	fmt.Println("  envguard generate-example")
	fmt.Println("  envguard sync-example")
}

func printValidateHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard validate")
	fmt.Println("  envguard validate --all")
	fmt.Println("  envguard validate --file .env.prod")
	fmt.Println("  envguard validate --file .env.prod --example .env.example.prod")
	fmt.Println("")
	fmt.Println("Validation checks:")
	fmt.Println("  - missing keys compared to the example file")
	fmt.Println("  - duplicate keys in the target env file")
	fmt.Println("  - unused keys not present in the example file")
	fmt.Println("  - typed values from examples/.env.types (optional)")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --all       Validate .env.dev, .env.prod, and .env.test")
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

func printScanCodeHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard scan-code")
	fmt.Println("  envguard scan-code --dir .")
	fmt.Println("  envguard scan-code --dir . --file .env.prod")
	fmt.Println("")
	fmt.Println("Checks:")
	fmt.Println("  - env variables used in code but missing in the env file")
	fmt.Println("  - env variables present in the env file but not used in code")
	fmt.Println("  - likely variable naming mismatches")
	fmt.Println("")
	fmt.Println("Supported patterns:")
	fmt.Println("  - os.Getenv(\"KEY\")")
	fmt.Println("  - os.LookupEnv(\"KEY\")")
	fmt.Println("  - process.env.KEY")
	fmt.Println("  - process.env[\"KEY\"]")
	fmt.Println("  - import.meta.env.KEY")
	fmt.Println("  - os.getenv(\"KEY\")")
	fmt.Println("  - os.environ[\"KEY\"]")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --dir       Root directory to scan")
	fmt.Println("  --file      Env file to compare against")
}

func printSecurityHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard security")
	fmt.Println("  envguard security --dir .")
	fmt.Println("  envguard security --file .env.prod")
	fmt.Println("  envguard security --dir . --file .env.prod")
	fmt.Println("")
	fmt.Println("Checks:")
	fmt.Println("  - secret-looking values in the env file")
	fmt.Println("  - secret-looking values in repository files")
	fmt.Println("  - secret-looking values in git history")
	fmt.Println("  - tracked env files in git")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --dir       Root directory to scan")
	fmt.Println("  --file      Env file to inspect")
}

func printLogScanHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard log-scan")
	fmt.Println("  envguard log-scan --dir .")
	fmt.Println("")
	fmt.Println("Checks:")
	fmt.Println("  - source code that logs env variable values")
	fmt.Println("  - log files containing secret-looking values")
	fmt.Println("  - log files containing sensitive key/value pairs")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --dir       Root directory to scan")
}

func printEncryptHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard encrypt")
	fmt.Println("  envguard encrypt --file .env.prod --out .env.prod.enc")
	fmt.Println("")
	fmt.Println("Encryption:")
	fmt.Println("  - reads the key from ENVGUARD_KEY")
	fmt.Println("  - encrypts the target env file using AES-GCM")
	fmt.Println("  - writes encrypted output to the output file")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Env file to encrypt")
	fmt.Println("  --out       Encrypted output file")
}

func printDecryptHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard decrypt")
	fmt.Println("  envguard decrypt --file .env.prod.enc --out .env.prod")
	fmt.Println("")
	fmt.Println("Decryption:")
	fmt.Println("  - reads the key from ENVGUARD_KEY")
	fmt.Println("  - decrypts an Env Guardian encrypted file")
	fmt.Println("  - writes plaintext output to the output file")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Encrypted file to decrypt")
	fmt.Println("  --out       Decrypted output file")
}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func hasAllFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--all" {
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
	case "scan-code":
		printScanCodeHelp()
		return 0
	case "security":
		printSecurityHelp()
		return 0
	case "log-scan":
		printLogScanHelp()
		return 0
	case "encrypt":
		printEncryptHelp()
		return 0
	case "decrypt":
		printDecryptHelp()
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
	allFlagSeen := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--all" {
			if allFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --all")
			}
			if fileFlagSeen || exampleFlagSeen {
				return "", "", fmt.Errorf("--all cannot be used with --file or --example")
			}
			allFlagSeen = true
			continue
		}

		if arg == "--file" {
			if allFlagSeen {
				return "", "", fmt.Errorf("--all cannot be used with --file or --example")
			}
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
			if allFlagSeen {
				return "", "", fmt.Errorf("--all cannot be used with --file or --example")
			}
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
			if allFlagSeen {
				return "", "", fmt.Errorf("--all cannot be used with --file or --example")
			}
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
			if allFlagSeen {
				return "", "", fmt.Errorf("--all cannot be used with --file or --example")
			}
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

func getScanCodeOptions(args []string) (string, string, error) {
	dirPath := "."
	envPath := ".env"
	dirFlagSeen := false
	fileFlagSeen := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--dir" {
			if dirFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --dir")
			}
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("missing value for --dir")
			}

			next := args[i+1]
			if strings.HasPrefix(next, "--") {
				return "", "", fmt.Errorf("missing value for --dir")
			}

			dirPath = next
			dirFlagSeen = true
			i++
			continue
		}

		if strings.HasPrefix(arg, "--dir=") {
			if dirFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --dir")
			}

			value := strings.TrimPrefix(arg, "--dir=")
			if value == "" {
				return "", "", fmt.Errorf("missing value for --dir")
			}

			dirPath = value
			dirFlagSeen = true
			continue
		}

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

		if strings.HasPrefix(arg, "--") {
			return "", "", fmt.Errorf("unknown flag: %s", arg)
		}

		return "", "", fmt.Errorf("unexpected argument: %s", arg)
	}

	return dirPath, envPath, nil
}

func getLogScanDir(args []string) (string, error) {
	dirPath := "."
	dirFlagSeen := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--dir" {
			if dirFlagSeen {
				return "", fmt.Errorf("duplicate flag: --dir")
			}
			if i+1 >= len(args) {
				return "", fmt.Errorf("missing value for --dir")
			}

			next := args[i+1]
			if strings.HasPrefix(next, "--") {
				return "", fmt.Errorf("missing value for --dir")
			}

			dirPath = next
			dirFlagSeen = true
			i++
			continue
		}

		if strings.HasPrefix(arg, "--dir=") {
			if dirFlagSeen {
				return "", fmt.Errorf("duplicate flag: --dir")
			}

			value := strings.TrimPrefix(arg, "--dir=")
			if value == "" {
				return "", fmt.Errorf("missing value for --dir")
			}

			dirPath = value
			dirFlagSeen = true
			continue
		}

		if strings.HasPrefix(arg, "--") {
			return "", fmt.Errorf("unknown flag: %s", arg)
		}

		return "", fmt.Errorf("unexpected argument: %s", arg)
	}

	return dirPath, nil
}

func getCryptoPaths(args []string, defaultInputPath string, defaultOutputPath string) (string, string, error) {
	inputPath := defaultInputPath
	outputPath := defaultOutputPath
	fileFlagSeen := false
	outFlagSeen := false

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

			inputPath = next
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

			inputPath = value
			fileFlagSeen = true
			continue
		}

		if arg == "--out" {
			if outFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --out")
			}
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("missing value for --out")
			}

			next := args[i+1]
			if strings.HasPrefix(next, "--") {
				return "", "", fmt.Errorf("missing value for --out")
			}

			outputPath = next
			outFlagSeen = true
			i++
			continue
		}

		if strings.HasPrefix(arg, "--out=") {
			if outFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --out")
			}

			value := strings.TrimPrefix(arg, "--out=")
			if value == "" {
				return "", "", fmt.Errorf("missing value for --out")
			}

			outputPath = value
			outFlagSeen = true
			continue
		}

		if strings.HasPrefix(arg, "--") {
			return "", "", fmt.Errorf("unknown flag: %s", arg)
		}

		return "", "", fmt.Errorf("unexpected argument: %s", arg)
	}

	if inputPath == outputPath {
		return "", "", fmt.Errorf("--file and --out must be different")
	}

	return inputPath, outputPath, nil
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

	for _, key := range result.UnusedKeys {
		fmt.Printf("[WARN] Unused key: %s\n", key)
		warningCount++
	}

	if errorCount == 0 && warningCount == 0 {
		fmt.Println("[PASS] Environment configuration looks good")
		fmt.Println("")
		fmt.Printf("Summary: %d error(s), %d warning(s)\n", errorCount, warningCount)
		return 0
	}

	if errorCount == 0 && warningCount > 0 {
		fmt.Println("")
		fmt.Println("[PASS] Environment configuration is valid with warnings")
		fmt.Println("")
		fmt.Printf("Summary: %d error(s), %d warning(s)\n", errorCount, warningCount)
		return 0
	}

	fmt.Println("")
	fmt.Printf("Summary: %d error(s), %d warning(s)\n", errorCount, warningCount)

	return 1
}

func runValidateAll() int {
	envTargets := []struct {
		envPath     string
		examplePath string
		label       string
	}{
		{envPath: ".env.dev", examplePath: ".env.example.dev", label: "dev"},
		{envPath: ".env.prod", examplePath: ".env.example.prod", label: "prod"},
		{envPath: ".env.test", examplePath: ".env.example.test", label: "test"},
	}

	finalExitCode := 0
	envFiles := make(map[string]*models.EnvFile)
	validationMissing := make(map[string]map[string]bool)

	fmt.Println("Env Validation Report")
	fmt.Println("---------------------")
	fmt.Println("Mode: --all")
	fmt.Println("")

	for i, target := range envTargets {
		fmt.Printf("[%s]\n", target.label)

		if _, err := os.Stat(target.envPath); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("[SKIP] %s missing\n", target.envPath)
			} else {
				fmt.Printf("[ERROR] could not access %s\n", target.envPath)
				finalExitCode = 1
			}

			if i < len(envTargets)-1 {
				fmt.Println("")
			}
			continue
		}

		if _, err := os.Stat(target.examplePath); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("[SKIP] %s missing\n", target.examplePath)
			} else {
				fmt.Printf("[ERROR] could not access %s\n", target.examplePath)
				finalExitCode = 1
			}

			if i < len(envTargets)-1 {
				fmt.Println("")
			}
			continue
		}

		envFile, err := parser.ParseEnvFile(target.envPath)
		exampleFile, err2 := parser.ParseEnvFile(target.examplePath)

		if err == nil && err2 == nil {
			envFiles[target.label] = envFile

			result := validator.ValidateEnv(envFile, exampleFile, map[string]string{})
			missingSet := make(map[string]bool)
			for _, key := range result.MissingKeys {
				missingSet[key] = true
			}
			validationMissing[target.label] = missingSet
		}

		exitCode := runValidate(target.envPath, target.examplePath)
		if exitCode != 0 {
			finalExitCode = 1
		}

		if i < len(envTargets)-1 {
			fmt.Println("")
		}
	}

	if len(envFiles) > 1 {
		fmt.Println("====================================")
		fmt.Println("Cross-Environment Consistency Check")
		fmt.Println("====================================")

		inconsistencies := validator.CompareEnvs(envFiles)

		printedConsistencyIssue := false

		if len(inconsistencies) == 0 {
			fmt.Println("[PASS] All environments are consistent")
		} else {
			for env, missingKeys := range inconsistencies {
				filtered := []string{}

				for _, key := range missingKeys {
					if validationMissing[env][key] {
						continue
					}
					filtered = append(filtered, key)
				}

				if len(filtered) == 0 {
					continue
				}

				printedConsistencyIssue = true

				fmt.Printf("[%s]\n", env)
				for _, key := range filtered {
					fmt.Printf("[WARNING] Missing key across environments: %s\n", key)
				}
				fmt.Println("")
			}

			if !printedConsistencyIssue {
				fmt.Println("[PASS] No additional cross-environment inconsistencies found")
			}
		}
	}

	return finalExitCode
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
	envTracked := security.IsGitTracked(".", envPath)
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

	securityWarningCount := 0
	if envTracked {
		fmt.Printf("\n[WARNING] %s appears to be tracked by git\n", envPath)
		securityWarningCount++
	}

	if len(missingInEnv) == 0 && envExists && exampleExists && securityWarningCount == 0 {
		fmt.Println("[PASS] Environment doctor checks passed")
		fmt.Println("")
		fmt.Printf("Summary: %d target file issue(s), %d example file issue(s), %d missing key(s), %d security warning(s)\n", targetFileIssues, exampleFileIssues, len(missingInEnv), securityWarningCount)
		return 0
	}

	fmt.Println("")
	fmt.Printf("Summary: %d target file issue(s), %d example file issue(s), %d missing key(s), %d security warning(s)\n", targetFileIssues, exampleFileIssues, len(missingInEnv), securityWarningCount)

	if len(missingInEnv) > 0 || !envExists {
		return 1
	}

	return 0
}

func runSecurity(dirPath string, envPath string) int {
	result, err := security.Run(dirPath, envPath)
	if err != nil {
		fmt.Printf("Error: could not run security checks: %v\n", err)
		return 1
	}

	fmt.Println("Env Security Report")
	fmt.Println("-------------------")
	fmt.Printf("Root directory: %s\n", dirPath)
	fmt.Printf("Env file: %s\n\n", envPath)

	secretFindingCount := 0
	warningCount := 0

	if result.EnvFileTracked {
		fmt.Printf("[WARNING] %s appears to be tracked by git\n", envPath)
		warningCount++
	}

	for _, finding := range result.EnvFindings {
		fmt.Printf("[ERROR] Env secret: %s (%s)\n", finding.Location, finding.Kind)
		secretFindingCount++
	}

	for _, finding := range result.RepositoryFindings {
		fmt.Printf("[ERROR] Repository secret: %s (%s)\n", finding.Location, finding.Kind)
		secretFindingCount++
	}

	for _, finding := range result.HistoryFindings {
		fmt.Printf("[ERROR] Git history secret: %s (%s)\n", finding.Location, finding.Kind)
		secretFindingCount++
	}

	if !result.HistoryScanned {
		fmt.Println("[WARNING] Git history scan skipped")
		warningCount++
	}

	if secretFindingCount == 0 && warningCount == 0 {
		fmt.Println("[PASS] Security checks found no issues")
		fmt.Println("")
		fmt.Printf("Summary: %d secret finding(s), %d warning(s)\n", secretFindingCount, warningCount)
		return 0
	}

	if secretFindingCount == 0 {
		fmt.Println("")
		fmt.Println("[PASS] Security checks completed with warnings")
		fmt.Println("")
		fmt.Printf("Summary: %d secret finding(s), %d warning(s)\n", secretFindingCount, warningCount)
		return 0
	}

	fmt.Println("")
	fmt.Printf("Summary: %d secret finding(s), %d warning(s)\n", secretFindingCount, warningCount)

	return 1
}

func runLogScan(dirPath string) int {
	result, err := logscan.Run(dirPath)
	if err != nil {
		fmt.Printf("Error: could not scan logs in %s\n", dirPath)
		return 1
	}

	sort.Slice(result.Findings, func(i, j int) bool {
		if result.Findings[i].Location == result.Findings[j].Location {
			return result.Findings[i].Kind < result.Findings[j].Kind
		}
		return result.Findings[i].Location < result.Findings[j].Location
	})

	fmt.Println("Env Log Exposure Report")
	fmt.Println("-----------------------")
	fmt.Printf("Root directory: %s\n", dirPath)
	fmt.Printf("Scanned files: %d\n\n", result.ScannedFilesCount)

	if len(result.Findings) == 0 {
		fmt.Println("[PASS] Log exposure scan found no issues")
		fmt.Println("")
		fmt.Printf("Summary: %d log exposure finding(s)\n", len(result.Findings))
		return 0
	}

	for _, finding := range result.Findings {
		fmt.Printf("[ERROR] Log exposure: %s (%s)\n", finding.Location, finding.Kind)
	}

	fmt.Println("")
	fmt.Printf("Summary: %d log exposure finding(s)\n", len(result.Findings))

	return 1
}

func runEncrypt(inputPath string, outputPath string) int {
	err := encryption.EncryptFile(inputPath, outputPath, os.Getenv("ENVGUARD_KEY"))
	if err != nil {
		fmt.Printf("Error: could not encrypt %s: %v\n", inputPath, err)
		return 1
	}

	fmt.Println("Env Encryption Report")
	fmt.Println("---------------------")
	fmt.Printf("Input file: %s\n", inputPath)
	fmt.Printf("Output file: %s\n\n", outputPath)
	fmt.Println("[PASS] Env file encrypted")

	return 0
}

func runDecrypt(inputPath string, outputPath string) int {
	err := encryption.DecryptFile(inputPath, outputPath, os.Getenv("ENVGUARD_KEY"))
	if err != nil {
		fmt.Printf("Error: could not decrypt %s: %v\n", inputPath, err)
		return 1
	}

	fmt.Println("Env Decryption Report")
	fmt.Println("---------------------")
	fmt.Printf("Input file: %s\n", inputPath)
	fmt.Printf("Output file: %s\n\n", outputPath)
	fmt.Println("[PASS] Env file decrypted")

	return 0
}

func runScanCode(dirPath string, envPath string) int {
	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		fmt.Printf("Error: could not read %s\n", envPath)
		return 1
	}

	result, err := codebase.Run(dirPath, envFile.Values)
	if err != nil {
		fmt.Printf("Error: could not scan codebase in %s\n", dirPath)
		return 1
	}

	fmt.Println("Codebase Analysis Report")
	fmt.Println("------------------------")
	fmt.Printf("Root directory: %s\n", dirPath)
	fmt.Printf("Env file: %s\n", envPath)
	fmt.Printf("Scanned files: %d\n", result.ScannedFilesCount)
	fmt.Printf("Detected env usages: %d\n", len(result.UsedKeys))

	if len(result.NamingMismatches) > 0 {
		fmt.Println("\nLikely naming mismatches:")
		for _, mismatch := range result.NamingMismatches {
			fmt.Printf("[WARN] Code uses %s but env file contains %s\n", mismatch.CodeKey, mismatch.EnvKey)
		}
	}

	if len(result.MissingInEnv) > 0 {
		fmt.Println("\nUsed in code but missing in env file:")
		for _, key := range result.MissingInEnv {
			fmt.Printf("[ERROR] Missing env key for code usage: %s\n", key)
		}
	}

	if len(result.UnusedInEnv) > 0 {
		fmt.Println("\nPresent in env file but unused in code:")
		for _, key := range result.UnusedInEnv {
			fmt.Printf("[WARN] Unused env key in codebase: %s\n", key)
		}
	}

	if len(result.NamingMismatches) == 0 && len(result.MissingInEnv) == 0 && len(result.UnusedInEnv) == 0 {
		fmt.Println("\n[PASS] Codebase analysis found no issues")
	}

	fmt.Println("")
	fmt.Printf(
		"Summary: %d naming mismatch(es), %d missing env key(s), %d unused env key(s)\n",
		len(result.NamingMismatches),
		len(result.MissingInEnv),
		len(result.UnusedInEnv),
	)

	if len(result.MissingInEnv) > 0 {
		return 1
	}

	return 0
}

func runGenerateExample(envPath string) int {
	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		fmt.Printf("Error: could not read %s\n", envPath)
		return 1
	}

	outputPath := ".env.example"

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error: could not create %s\n", outputPath)
		return 1
	}
	defer file.Close()

	keys := make([]string, 0, len(envFile.Values))
	for key := range envFile.Values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		_, err := file.WriteString(fmt.Sprintf("%s=\n", key))
		if err != nil {
			fmt.Println("Error: failed to write to file")
			return 1
		}
	}

	fmt.Println("Env Example Generated")
	fmt.Println("---------------------")
	fmt.Printf("Source file: %s\n", envPath)
	fmt.Printf("Output file: %s\n\n", outputPath)
	fmt.Printf("Generated %d keys\n", len(keys))

	return 0
}

func runSyncExample(envPath string) int {
	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		fmt.Printf("Error: could not read %s\n", envPath)
		return 1
	}

	examplePath := ".env.example"

	exampleFile, err := parser.ParseEnvFile(examplePath)

	existingKeys := map[string]bool{}
	if err == nil {
		for key := range exampleFile.Values {
			existingKeys[key] = true
		}
	}

	file, err := os.OpenFile(examplePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error: could not open %s\n", examplePath)
		return 1
	}
	defer file.Close()

	newKeys := []string{}

	for key := range envFile.Values {
		if !existingKeys[key] {
			newKeys = append(newKeys, key)
		}
	}

	sort.Strings(newKeys)

	for _, key := range newKeys {
		_, err := file.WriteString(fmt.Sprintf("%s=\n", key))
		if err != nil {
			fmt.Println("Error: failed to write to file")
			return 1
		}
	}

	fmt.Println("Env Example Sync Report")
	fmt.Println("------------------------")
	fmt.Printf("Source file: %s\n", envPath)
	fmt.Printf("Synced file: %s\n\n", examplePath)

	if len(newKeys) == 0 {
		fmt.Println("[PASS] .env.example already up to date")
	} else {
		fmt.Printf("Added %d new key(s):\n", len(newKeys))
		for _, key := range newKeys {
			fmt.Printf("- %s\n", key)
		}
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
		if hasAllFlag(args[1:]) {
			envPath, examplePath, err := getValidatePaths(args[1:])
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			if envPath != ".env" || examplePath != ".env.example" {
				fmt.Println("Error: --all cannot be used with --file or --example")
				os.Exit(1)
			}
			os.Exit(runValidateAll())
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
	case "scan-code":
		if hasHelpFlag(args[1:]) {
			printScanCodeHelp()
			os.Exit(0)
		}
		dirPath, envPath, err := getScanCodeOptions(args[1:])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runScanCode(dirPath, envPath))
	case "security":
		if hasHelpFlag(args[1:]) {
			printSecurityHelp()
			os.Exit(0)
		}
		dirPath, envPath, err := getScanCodeOptions(args[1:])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runSecurity(dirPath, envPath))
	case "log-scan":
		if hasHelpFlag(args[1:]) {
			printLogScanHelp()
			os.Exit(0)
		}
		dirPath, err := getLogScanDir(args[1:])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runLogScan(dirPath))
	case "encrypt":
		if hasHelpFlag(args[1:]) {
			printEncryptHelp()
			os.Exit(0)
		}
		inputPath, outputPath, err := getCryptoPaths(args[1:], ".env", ".env.enc")
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runEncrypt(inputPath, outputPath))
	case "decrypt":
		if hasHelpFlag(args[1:]) {
			printDecryptHelp()
			os.Exit(0)
		}
		inputPath, outputPath, err := getCryptoPaths(args[1:], ".env.enc", ".env.decrypted")
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runDecrypt(inputPath, outputPath))
	case "generate-example":
		envPath, err := getLintFilePath(args[1:])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runGenerateExample(envPath))
	case "sync-example":
		envPath, err := getLintFilePath(args[1:])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runSyncExample(envPath))
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
