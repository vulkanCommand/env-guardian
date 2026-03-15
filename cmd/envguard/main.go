package main

import (
	"fmt"
	"os"

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
		fmt.Println("validate command not implemented yet")
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