package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/vulkanCommand/env-guardian/internal/analyzer"
	"github.com/vulkanCommand/env-guardian/internal/doctor"
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
	fmt.Println("  envguard version")
	fmt.Println("  envguard validate")
	fmt.Println("  envguard lint")
	fmt.Println("  envguard analyze")
	fmt.Println("  envguard doctor")
}

func runValidate() int {
	envFile, err := parser.ParseEnvFile(".env")
	if err != nil {
		fmt.Println("Error: could not read .env")
		return 1
	}

	exampleFile, err := parser.ParseEnvFile(".env.example")
	if err != nil {
		fmt.Println("Error: could not read .env.example")
		return 1
	}

	result := validator.ValidateEnv(envFile, exampleFile)

	sort.Strings(result.MissingKeys)
	sort.Strings(result.DuplicateKeys)

	fmt.Println("Env Validation Report")
	fmt.Println("---------------------")

	errorCount := 0

	for _, key := range result.MissingKeys {
		fmt.Printf("[ERROR] Missing key: %s\n", key)
		errorCount++
	}

	for _, key := range result.DuplicateKeys {
		fmt.Printf("[ERROR] Duplicate key: %s\n", key)
		errorCount++
	}

	if errorCount == 0 {
		fmt.Println("[PASS] Environment configuration looks good")
		return 0
	}

	fmt.Println("")
	fmt.Printf("Summary: %d error(s) found\n", errorCount)

	return 1
}

func runLint() int {
	result, err := linter.Run(".env")
	if err != nil {
		fmt.Println("Error: could not read .env")
		return 1
	}

	fmt.Println("Env Lint Report")
	fmt.Println("---------------")

	if len(result.InvalidLines) == 0 {
		fmt.Println("[PASS] .env syntax looks good")
		return 0
	}

	for _, issue := range result.InvalidLines {
		fmt.Printf("[ERROR] %s\n", issue)
	}

	fmt.Println("")
	fmt.Printf("Summary: %d lint issue(s) found\n", len(result.InvalidLines))

	return 1
}

func runAnalyze() int {
	result, err := analyzer.Run(".env")
	if err != nil {
		fmt.Println("Error: could not read .env")
		return 1
	}

	fmt.Println("Env Analysis Report")
	fmt.Println("-------------------")
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

	return 0
}

func runDoctor() int {
	result := doctor.Run()

	fmt.Println("Env Doctor Report")
	fmt.Println("------------------")

	if result.EnvFileExists {
		fmt.Println("[OK] .env file exists")
	} else {
		fmt.Println("[ERROR] .env file missing")
	}

	if result.ExampleFileExists {
		fmt.Println("[OK] .env.example file exists")
	} else {
		fmt.Println("[WARNING] .env.example file missing")
	}

	if len(result.MissingInEnv) > 0 {
		fmt.Println("\n[WARNING] Missing keys in .env (present in .env.example):")
		for _, key := range result.MissingInEnv {
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
		os.Exit(runValidate())
	case "lint":
		os.Exit(runLint())
	case "analyze":
		os.Exit(runAnalyze())
	case "doctor":
		os.Exit(runDoctor())
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("unknown command: %s\n\n", args[0])
		printHelp()
		os.Exit(1)
	}
}