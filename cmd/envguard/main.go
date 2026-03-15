package main

import (
	"fmt"
	"os"
	"sort"

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
		fmt.Println("lint command not implemented yet")
	case "analyze":
		fmt.Println("analyze command not implemented yet")
	case "doctor":
		fmt.Println("doctor command not implemented yet")
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("unknown command: %s\n\n", args[0])
		printHelp()
		os.Exit(1)
	}
}