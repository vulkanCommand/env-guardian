package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/vulkanCommand/env-guardian/internal/analyzer"
	"github.com/vulkanCommand/env-guardian/internal/codebase"
	"github.com/vulkanCommand/env-guardian/internal/encryption"
	"github.com/vulkanCommand/env-guardian/internal/linter"
	"github.com/vulkanCommand/env-guardian/internal/logscan"
	"github.com/vulkanCommand/env-guardian/internal/models"
	"github.com/vulkanCommand/env-guardian/internal/parser"
	"github.com/vulkanCommand/env-guardian/internal/runtimecheck"
	"github.com/vulkanCommand/env-guardian/internal/security"
	"github.com/vulkanCommand/env-guardian/internal/validator"
	"github.com/vulkanCommand/env-guardian/internal/version"
)

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiGreen  = "\033[32m"
	ansiRed    = "\033[31m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
)

func colorEnabled() bool {
	return os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"
}

func colorText(value string, color string) string {
	if !colorEnabled() {
		return value
	}
	return color + value + ansiReset
}

func green(value string) string {
	return colorText(value, ansiGreen)
}

func red(value string) string {
	return colorText(value, ansiRed)
}

func yellow(value string) string {
	return colorText(value, ansiYellow)
}

func cyan(value string) string {
	return colorText(value, ansiCyan)
}

func bold(value string) string {
	return colorText(value, ansiBold)
}

func statusLine(value string) string {
	replacer := strings.NewReplacer(
		"[PASS]", green("[PASS]"),
		"[ERROR]", red("[ERROR]"),
		"[WARN]", yellow("[WARN]"),
		"[WARNING]", yellow("[WARNING]"),
		"[OK]", green("[OK]"),
		"[RUN]", cyan("[RUN]"),
		"[SKIP]", yellow("[SKIP]"),
		"Error:", red("Error:"),
	)

	return replacer.Replace(value)
}

func printStatusLine(value string) {
	fmt.Println(statusLine(value))
}

func printfStatusLine(format string, args ...any) {
	fmt.Printf(statusLine(format), args...)
}

func printHelp() {
	printTitleCard()
	fmt.Println("")
	fmt.Println(green("COMMANDS"))
	fmt.Println(green("+-------------------+--------------------------------------------------+"))
	fmt.Println(green("| Command           | Description                                      |"))
	fmt.Println(green("+-------------------+--------------------------------------------------+"))
	fmt.Println(green("| validate          | Validate .env against .env.example               |"))
	fmt.Println(green("| lint              | Check env file syntax                            |"))
	fmt.Println(green("| analyze           | Inspect empty values and sensitive-looking keys  |"))
	fmt.Println(green("| doctor            | Diagnose files, keys, and tracked env files      |"))
	fmt.Println(green("| scan-code         | Compare env usage in code against an env file    |"))
	fmt.Println(green("| security          | Scan env, repository files, and git history      |"))
	fmt.Println(green("| log-scan          | Detect accidental env value logging              |"))
	fmt.Println(green("| encrypt           | Encrypt an env file with ENVGUARD_KEY            |"))
	fmt.Println(green("| decrypt           | Decrypt an Env Guardian encrypted file           |"))
	fmt.Println(green("| docker            | Validate Dockerfile env references               |"))
	fmt.Println(green("| ci                | Run fail-fast validation for CI                  |"))
	fmt.Println(green("| run               | Validate env config before starting a command    |"))
	fmt.Println(green("| generate-example  | Create .env.example from an env file             |"))
	fmt.Println(green("| sync-example      | Append missing keys to .env.example              |"))
	fmt.Println(green("| version           | Print the Env Guardian version                   |"))
	fmt.Println(green("+-------------------+--------------------------------------------------+"))
	fmt.Println("")
	fmt.Println(green("QUICK START"))
	fmt.Println(cyan("  envguard validate"))
	fmt.Println(cyan("  envguard security"))
	fmt.Println(cyan("  envguard ci --json"))
	fmt.Println("")
	fmt.Println(green("HELP"))
	fmt.Println(cyan("  envguard help <command>"))
	fmt.Println(cyan("  envguard help validate"))
	fmt.Println("")
	printSupport()
}

func printTitleCard() {
	fmt.Println(green("=========================================================================="))
	fmt.Println(green("  ______ _   ___     __    _____ _    _          _____  _____ _____          _   _ "))
	fmt.Println(green(" |  ____| \\ | \\ \\   / /   / ____| |  | |   /\\   |  __ \\|  __ \\_   _|   /\\   | \\ | |"))
	fmt.Println(green(" | |__  |  \\| |\\ \\_/ /   | |  __| |  | |  /  \\  | |__) | |  | || |    /  \\  |  \\| |"))
	fmt.Println(green(" |  __| | . ` | \\   /    | | |_ | |  | | / /\\ \\ |  _  /| |  | || |   / /\\ \\ | . ` |"))
	fmt.Println(green(" | |____| |\\  |  | |     | |__| | |__| |/ ____ \\| | \\ \\| |__| || |_ / ____ \\| |\\  |"))
	fmt.Println(green(" |______|_| \\_|  |_|      \\_____|\\____//_/    \\_\\_|  \\_\\_____/_____/_/    \\_\\_| \\_|"))
	fmt.Println(green("=========================================================================="))
	fmt.Println(green("           ENV GUARDIAN CLI") + "  " + bold("(v"+version.Version+")"))
	fmt.Println(green("           Validate. Secure. Encrypt. Ship environment files safely."))
	fmt.Println(green("=========================================================================="))
}

func printSupport() {
	fmt.Println(green("SUPPORT"))
	fmt.Println(cyan("  Email:  gdkalyan2109@gmail.com"))
	fmt.Println(cyan("  Issues: https://github.com/vulkanCommand/env-guardian/issues"))
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
	fmt.Println("  --json      Print machine-readable JSON output")
}

func printLintHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard lint")
	fmt.Println("  envguard lint --file .env.prod")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Target env file to lint")
	fmt.Println("  --json      Print machine-readable JSON output")
}

func printAnalyzeHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard analyze")
	fmt.Println("  envguard analyze --file .env.prod")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Target env file to analyze")
	fmt.Println("  --json      Print machine-readable JSON output")
}

func printDoctorHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard doctor")
	fmt.Println("  envguard doctor --file .env.prod --example .env.example.prod")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Target env file to inspect")
	fmt.Println("  --example   Example env file to compare against")
	fmt.Println("  --json      Print machine-readable JSON output")
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
	fmt.Println("  --json      Print machine-readable JSON output")
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
	fmt.Println("  --json      Print machine-readable JSON output")
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
	fmt.Println("  --json      Print machine-readable JSON output")
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

func printDockerHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard docker")
	fmt.Println("  envguard docker --dockerfile Dockerfile --file .env.prod")
	fmt.Println("")
	fmt.Println("Checks:")
	fmt.Println("  - Dockerfile ARG and ENV variables")
	fmt.Println("  - Dockerfile $KEY and ${KEY} references")
	fmt.Println("  - missing Docker runtime keys in the env file")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --dockerfile   Dockerfile path to inspect")
	fmt.Println("  --file         Env file to compare against")
	fmt.Println("  --json         Print machine-readable JSON output")
}

func printCIHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard ci")
	fmt.Println("  envguard ci --file .env.prod --example .env.example.prod")
	fmt.Println("")
	fmt.Println("Checks:")
	fmt.Println("  - env syntax linting")
	fmt.Println("  - required key validation")
	fmt.Println("  - duplicate key validation")
	fmt.Println("  - optional type validation")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Target env file to validate")
	fmt.Println("  --example   Example env file to compare against")
	fmt.Println("  --json      Print machine-readable JSON output")
}

func printRunHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard run -- <command>")
	fmt.Println("  envguard run --file .env.prod --example .env.example.prod -- <command>")
	fmt.Println("")
	fmt.Println("Runtime wrapper:")
	fmt.Println("  - validates env configuration before starting a command")
	fmt.Println("  - runs the command only when validation passes")
	fmt.Println("  - returns the wrapped command exit code")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Target env file to validate")
	fmt.Println("  --example   Example env file to compare against")
}

func printGenerateExampleHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard generate-example")
	fmt.Println("  envguard generate-example --file .env.prod")
	fmt.Println("")
	fmt.Println("Workflow:")
	fmt.Println("  - reads keys from the target env file")
	fmt.Println("  - writes .env.example with empty values")
	fmt.Println("  - overwrites the existing .env.example file")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Source env file to read")
}

func printSyncExampleHelp() {
	fmt.Println("Usage:")
	fmt.Println("  envguard sync-example")
	fmt.Println("  envguard sync-example --file .env.prod")
	fmt.Println("")
	fmt.Println("Workflow:")
	fmt.Println("  - reads keys from the target env file")
	fmt.Println("  - appends missing keys to .env.example")
	fmt.Println("  - does not overwrite existing example keys")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --file      Source env file to read")
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

func extractJSONFlag(args []string) ([]string, bool, error) {
	filtered := []string{}
	jsonSeen := false

	for _, arg := range args {
		if arg == "--json" {
			if jsonSeen {
				return nil, false, fmt.Errorf("duplicate flag: --json")
			}
			jsonSeen = true
			continue
		}

		filtered = append(filtered, arg)
	}

	return filtered, jsonSeen, nil
}

func handleHelpCommand(args []string) int {
	if len(args) == 0 {
		printHelp()
		return 0
	}

	if len(args) > 1 {
		printfStatusLine("Error: unexpected argument: %s\n", args[1])
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
	case "docker":
		printDockerHelp()
		return 0
	case "ci":
		printCIHelp()
		return 0
	case "run":
		printRunHelp()
		return 0
	case "generate-example":
		printGenerateExampleHelp()
		return 0
	case "sync-example":
		printSyncExampleHelp()
		return 0
	default:
		printfStatusLine("Error: unknown help topic: %s\n", args[0])
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

func getDockerOptions(args []string) (string, string, error) {
	dockerfilePath := "Dockerfile"
	envPath := ".env"
	dockerfileFlagSeen := false
	fileFlagSeen := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--dockerfile" {
			if dockerfileFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --dockerfile")
			}
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("missing value for --dockerfile")
			}

			next := args[i+1]
			if strings.HasPrefix(next, "--") {
				return "", "", fmt.Errorf("missing value for --dockerfile")
			}

			dockerfilePath = next
			dockerfileFlagSeen = true
			i++
			continue
		}

		if strings.HasPrefix(arg, "--dockerfile=") {
			if dockerfileFlagSeen {
				return "", "", fmt.Errorf("duplicate flag: --dockerfile")
			}

			value := strings.TrimPrefix(arg, "--dockerfile=")
			if value == "" {
				return "", "", fmt.Errorf("missing value for --dockerfile")
			}

			dockerfilePath = value
			dockerfileFlagSeen = true
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

	return dockerfilePath, envPath, nil
}

func getRunOptions(args []string) (string, string, []string, error) {
	separatorIndex := -1
	for i, arg := range args {
		if arg == "--" {
			separatorIndex = i
			break
		}
	}

	if separatorIndex == -1 {
		return "", "", nil, fmt.Errorf("missing command separator: --")
	}

	envPath, examplePath, err := getValidatePaths(args[:separatorIndex])
	if err != nil {
		return "", "", nil, err
	}

	commandArgs := args[separatorIndex+1:]
	if len(commandArgs) == 0 {
		return "", "", nil, fmt.Errorf("missing command to run")
	}

	return envPath, examplePath, commandArgs, nil
}

type jsonReport struct {
	Command       string         `json:"command"`
	Status        string         `json:"status"`
	TargetFile    string         `json:"target_file,omitempty"`
	ExampleFile   string         `json:"example_file,omitempty"`
	RootDirectory string         `json:"root_directory,omitempty"`
	Dockerfile    string         `json:"dockerfile,omitempty"`
	Errors        []string       `json:"errors"`
	Warnings      []string       `json:"warnings"`
	Summary       map[string]int `json:"summary"`
	Details       any            `json:"details,omitempty"`
}

func newJSONReport(command string) jsonReport {
	return jsonReport{
		Command:  command,
		Status:   "pass",
		Errors:   []string{},
		Warnings: []string{},
		Summary:  map[string]int{},
	}
}

func printJSONReport(report jsonReport) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(report)
}

func reportStatus(errorCount int, warningCount int) string {
	if errorCount > 0 {
		return "fail"
	}
	if warningCount > 0 {
		return "warning"
	}
	return "pass"
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func loadOptionalTypeSchema() (map[string]string, error) {
	schema := map[string]string{}

	loadedSchema, err := parser.LoadTypeSchema("examples/.env.types")
	if err != nil {
		if os.IsNotExist(err) {
			return schema, nil
		}
		return schema, err
	}

	return loadedSchema, nil
}

func runValidateJSON(envPath string, examplePath string) int {
	report := newJSONReport("validate")
	report.TargetFile = envPath
	report.ExampleFile = examplePath

	schema, err := loadOptionalTypeSchema()
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read type schema file: %v", err))
		report.Summary["errors"] = 1
		report.Summary["warnings"] = 0
		printJSONReport(report)
		return 1
	}

	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read %s", envPath))
		report.Summary["errors"] = 1
		report.Summary["warnings"] = 0
		printJSONReport(report)
		return 1
	}

	exampleFile, err := parser.ParseEnvFile(examplePath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read %s", examplePath))
		report.Summary["errors"] = 1
		report.Summary["warnings"] = 0
		printJSONReport(report)
		return 1
	}

	result := validator.ValidateEnv(envFile, exampleFile, schema)

	sort.Strings(result.MissingKeys)
	sort.Strings(result.DuplicateKeys)
	sort.Strings(result.UnusedKeys)
	sort.Strings(result.InvalidTypeValues)

	for _, key := range result.MissingKeys {
		report.Errors = append(report.Errors, fmt.Sprintf("Missing key: %s", key))
	}
	for _, key := range result.DuplicateKeys {
		report.Errors = append(report.Errors, fmt.Sprintf("Duplicate key: %s", key))
	}
	for _, issue := range result.InvalidTypeValues {
		report.Errors = append(report.Errors, fmt.Sprintf("Invalid type: %s", issue))
	}
	for _, key := range result.UnusedKeys {
		report.Warnings = append(report.Warnings, fmt.Sprintf("Unused key: %s", key))
	}

	errorCount := len(report.Errors)
	warningCount := len(report.Warnings)
	report.Status = reportStatus(errorCount, warningCount)
	report.Summary["errors"] = errorCount
	report.Summary["warnings"] = warningCount
	report.Details = map[string]any{
		"missing_keys":        result.MissingKeys,
		"duplicate_keys":      result.DuplicateKeys,
		"unused_keys":         result.UnusedKeys,
		"invalid_type_values": result.InvalidTypeValues,
		"type_schema_loaded":  len(schema) > 0,
	}

	printJSONReport(report)

	if errorCount > 0 {
		return 1
	}
	return 0
}

func runValidateAllJSON() int {
	report := newJSONReport("validate")
	report.Summary["errors"] = 0
	report.Summary["warnings"] = 0

	schema, err := loadOptionalTypeSchema()
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read type schema file: %v", err))
		report.Summary["errors"] = 1
		report.Summary["warnings"] = 0
		printJSONReport(report)
		return 1
	}

	envTargets := []struct {
		envPath     string
		examplePath string
		label       string
	}{
		{envPath: ".env.dev", examplePath: ".env.example.dev", label: "dev"},
		{envPath: ".env.prod", examplePath: ".env.example.prod", label: "prod"},
		{envPath: ".env.test", examplePath: ".env.example.test", label: "test"},
	}

	type environmentReport struct {
		Label       string   `json:"label"`
		TargetFile  string   `json:"target_file"`
		ExampleFile string   `json:"example_file"`
		Status      string   `json:"status"`
		Skipped     bool     `json:"skipped"`
		Errors      []string `json:"errors"`
		Warnings    []string `json:"warnings"`
	}

	environments := []environmentReport{}
	envFiles := make(map[string]*models.EnvFile)
	validationMissing := make(map[string]map[string]bool)
	totalErrors := 0
	totalWarnings := 0

	for _, target := range envTargets {
		environment := environmentReport{
			Label:       target.label,
			TargetFile:  target.envPath,
			ExampleFile: target.examplePath,
			Status:      "pass",
			Errors:      []string{},
			Warnings:    []string{},
		}

		if _, err := os.Stat(target.envPath); err != nil {
			if os.IsNotExist(err) {
				environment.Skipped = true
				environment.Status = "skipped"
			} else {
				message := fmt.Sprintf("could not access %s", target.envPath)
				environment.Status = "fail"
				environment.Errors = append(environment.Errors, message)
				report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", target.label, message))
				totalErrors++
			}
			environments = append(environments, environment)
			continue
		}

		if _, err := os.Stat(target.examplePath); err != nil {
			if os.IsNotExist(err) {
				environment.Skipped = true
				environment.Status = "skipped"
			} else {
				message := fmt.Sprintf("could not access %s", target.examplePath)
				environment.Status = "fail"
				environment.Errors = append(environment.Errors, message)
				report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", target.label, message))
				totalErrors++
			}
			environments = append(environments, environment)
			continue
		}

		envFile, err := parser.ParseEnvFile(target.envPath)
		if err != nil {
			message := fmt.Sprintf("could not read %s", target.envPath)
			environment.Status = "fail"
			environment.Errors = append(environment.Errors, message)
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", target.label, message))
			totalErrors++
			environments = append(environments, environment)
			continue
		}

		exampleFile, err := parser.ParseEnvFile(target.examplePath)
		if err != nil {
			message := fmt.Sprintf("could not read %s", target.examplePath)
			environment.Status = "fail"
			environment.Errors = append(environment.Errors, message)
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", target.label, message))
			totalErrors++
			environments = append(environments, environment)
			continue
		}

		envFiles[target.label] = envFile
		result := validator.ValidateEnv(envFile, exampleFile, schema)

		sort.Strings(result.MissingKeys)
		sort.Strings(result.DuplicateKeys)
		sort.Strings(result.UnusedKeys)
		sort.Strings(result.InvalidTypeValues)

		missingSet := make(map[string]bool)
		for _, key := range result.MissingKeys {
			message := fmt.Sprintf("Missing key: %s", key)
			environment.Errors = append(environment.Errors, message)
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", target.label, message))
			missingSet[key] = true
		}
		validationMissing[target.label] = missingSet

		for _, key := range result.DuplicateKeys {
			message := fmt.Sprintf("Duplicate key: %s", key)
			environment.Errors = append(environment.Errors, message)
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", target.label, message))
		}
		for _, issue := range result.InvalidTypeValues {
			message := fmt.Sprintf("Invalid type: %s", issue)
			environment.Errors = append(environment.Errors, message)
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", target.label, message))
		}
		for _, key := range result.UnusedKeys {
			message := fmt.Sprintf("Unused key: %s", key)
			environment.Warnings = append(environment.Warnings, message)
			report.Warnings = append(report.Warnings, fmt.Sprintf("%s: %s", target.label, message))
		}

		totalErrors += len(environment.Errors)
		totalWarnings += len(environment.Warnings)
		environment.Status = reportStatus(len(environment.Errors), len(environment.Warnings))
		environments = append(environments, environment)
	}

	consistencyWarnings := []string{}
	if len(envFiles) > 1 {
		inconsistencies := validator.CompareEnvs(envFiles)
		envNames := make([]string, 0, len(inconsistencies))
		for env := range inconsistencies {
			envNames = append(envNames, env)
		}
		sort.Strings(envNames)

		for _, env := range envNames {
			missingKeys := inconsistencies[env]
			sort.Strings(missingKeys)
			for _, key := range missingKeys {
				if validationMissing[env][key] {
					continue
				}
				consistencyWarnings = append(consistencyWarnings, fmt.Sprintf("%s missing key across environments: %s", env, key))
			}
		}
	}

	report.Warnings = append(report.Warnings, consistencyWarnings...)
	totalWarnings += len(consistencyWarnings)
	report.Status = reportStatus(totalErrors, totalWarnings)
	report.Summary["errors"] = totalErrors
	report.Summary["warnings"] = totalWarnings
	report.Details = map[string]any{
		"mode":                 "all",
		"environments":         environments,
		"consistency_warnings": consistencyWarnings,
		"type_schema_loaded":   len(schema) > 0,
	}

	printJSONReport(report)

	if totalErrors > 0 {
		return 1
	}
	return 0
}

func runValidate(envPath string, examplePath string) int {
	schema := map[string]string{}

	loadedSchema, err := parser.LoadTypeSchema("examples/.env.types")
	if err != nil {
		if !os.IsNotExist(err) {
			printfStatusLine("Error: could not read type schema file: %v\n", err)
			return 1
		}
	} else {
		schema = loadedSchema
	}

	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		printfStatusLine("Error: could not read %s\n", envPath)
		return 1
	}

	exampleFile, err := parser.ParseEnvFile(examplePath)
	if err != nil {
		printfStatusLine("Error: could not read %s\n", examplePath)
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
		printfStatusLine("[ERROR] Missing key: %s\n", key)
		errorCount++
	}

	for _, key := range result.DuplicateKeys {
		printfStatusLine("[ERROR] Duplicate key: %s\n", key)
		errorCount++
	}

	for _, issue := range result.InvalidTypeValues {
		printfStatusLine("[ERROR] Invalid type: %s\n", issue)
		errorCount++
	}

	for _, key := range result.UnusedKeys {
		printfStatusLine("[WARN] Unused key: %s\n", key)
		warningCount++
	}

	if errorCount == 0 && warningCount == 0 {
		printStatusLine("[PASS] Environment configuration looks good")
		fmt.Println("")
		fmt.Printf("Summary: %d error(s), %d warning(s)\n", errorCount, warningCount)
		return 0
	}

	if errorCount == 0 && warningCount > 0 {
		fmt.Println("")
		printStatusLine("[PASS] Environment configuration is valid with warnings")
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
				printfStatusLine("[SKIP] %s missing\n", target.envPath)
			} else {
				printfStatusLine("[ERROR] could not access %s\n", target.envPath)
				finalExitCode = 1
			}

			if i < len(envTargets)-1 {
				fmt.Println("")
			}
			continue
		}

		if _, err := os.Stat(target.examplePath); err != nil {
			if os.IsNotExist(err) {
				printfStatusLine("[SKIP] %s missing\n", target.examplePath)
			} else {
				printfStatusLine("[ERROR] could not access %s\n", target.examplePath)
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
			printStatusLine("[PASS] All environments are consistent")
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
					printfStatusLine("[WARNING] Missing key across environments: %s\n", key)
				}
				fmt.Println("")
			}

			if !printedConsistencyIssue {
				printStatusLine("[PASS] No additional cross-environment inconsistencies found")
			}
		}
	}

	return finalExitCode
}

func runLint(envPath string) int {
	result, err := linter.Run(envPath)
	if err != nil {
		printfStatusLine("Error: could not read %s\n", envPath)
		return 1
	}

	fmt.Println("Env Lint Report")
	fmt.Println("---------------")
	fmt.Printf("Target file: %s\n\n", envPath)

	if len(result.InvalidLines) == 0 {
		printStatusLine("[PASS] Env syntax looks good")
		fmt.Println("")
		fmt.Printf("Summary: %d lint issue(s) found\n", len(result.InvalidLines))
		return 0
	}

	for _, issue := range result.InvalidLines {
		printfStatusLine("[ERROR] %s\n", issue)
	}

	fmt.Println("")
	fmt.Printf("Summary: %d lint issue(s) found\n", len(result.InvalidLines))

	return 1
}

func runAnalyze(envPath string) int {
	result, err := analyzer.Run(envPath)
	if err != nil {
		printfStatusLine("Error: could not read %s\n", envPath)
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
		printStatusLine("[PASS] Environment analysis found no issues")
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
		printfStatusLine("[OK] %s file exists\n", envPath)
	} else {
		printfStatusLine("[ERROR] %s file missing\n", envPath)
	}

	if exampleExists {
		printfStatusLine("[OK] %s file exists\n", examplePath)
	} else {
		printfStatusLine("[WARNING] %s file missing\n", examplePath)
	}

	if len(missingInEnv) > 0 {
		printfStatusLine("\n[WARNING] Missing keys in %s (present in %s):\n", envPath, examplePath)
		for _, key := range missingInEnv {
			fmt.Printf("- %s\n", key)
		}
	}

	securityWarningCount := 0
	if envTracked {
		printfStatusLine("\n[WARNING] %s appears to be tracked by git\n", envPath)
		securityWarningCount++
	}

	if len(missingInEnv) == 0 && envExists && exampleExists && securityWarningCount == 0 {
		printStatusLine("[PASS] Environment doctor checks passed")
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
		printfStatusLine("Error: could not run security checks: %v\n", err)
		return 1
	}

	fmt.Println("Env Security Report")
	fmt.Println("-------------------")
	fmt.Printf("Root directory: %s\n", dirPath)
	fmt.Printf("Env file: %s\n\n", envPath)

	secretFindingCount := 0
	warningCount := 0

	if result.EnvFileTracked {
		printfStatusLine("[WARNING] %s appears to be tracked by git\n", envPath)
		warningCount++
	}

	for _, finding := range result.EnvFindings {
		printfStatusLine("[ERROR] Env secret: %s (%s)\n", finding.Location, finding.Kind)
		secretFindingCount++
	}

	for _, finding := range result.RepositoryFindings {
		printfStatusLine("[ERROR] Repository secret: %s (%s)\n", finding.Location, finding.Kind)
		secretFindingCount++
	}

	for _, finding := range result.HistoryFindings {
		printfStatusLine("[ERROR] Git history secret: %s (%s)\n", finding.Location, finding.Kind)
		secretFindingCount++
	}

	if !result.HistoryScanned {
		printStatusLine("[WARNING] Git history scan skipped")
		warningCount++
	}

	if secretFindingCount == 0 && warningCount == 0 {
		printStatusLine("[PASS] Security checks found no issues")
		fmt.Println("")
		fmt.Printf("Summary: %d secret finding(s), %d warning(s)\n", secretFindingCount, warningCount)
		return 0
	}

	if secretFindingCount == 0 {
		fmt.Println("")
		printStatusLine("[PASS] Security checks completed with warnings")
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
		printfStatusLine("Error: could not scan logs in %s\n", dirPath)
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
		printStatusLine("[PASS] Log exposure scan found no issues")
		fmt.Println("")
		fmt.Printf("Summary: %d log exposure finding(s)\n", len(result.Findings))
		return 0
	}

	for _, finding := range result.Findings {
		printfStatusLine("[ERROR] Log exposure: %s (%s)\n", finding.Location, finding.Kind)
	}

	fmt.Println("")
	fmt.Printf("Summary: %d log exposure finding(s)\n", len(result.Findings))

	return 1
}

func runLintJSON(envPath string) int {
	report := newJSONReport("lint")
	report.TargetFile = envPath

	result, err := linter.Run(envPath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read %s", envPath))
		report.Summary["lint_issues"] = 1
		printJSONReport(report)
		return 1
	}

	result.InvalidLines = nonNilStrings(result.InvalidLines)
	sort.Strings(result.InvalidLines)
	report.Errors = append(report.Errors, result.InvalidLines...)
	report.Status = reportStatus(len(report.Errors), 0)
	report.Summary["lint_issues"] = len(result.InvalidLines)
	report.Details = map[string]any{
		"invalid_lines": result.InvalidLines,
	}

	printJSONReport(report)

	if len(result.InvalidLines) > 0 {
		return 1
	}
	return 0
}

func runAnalyzeJSON(envPath string) int {
	report := newJSONReport("analyze")
	report.TargetFile = envPath

	result, err := analyzer.Run(envPath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read %s", envPath))
		report.Summary["errors"] = 1
		report.Summary["warnings"] = 0
		printJSONReport(report)
		return 1
	}

	result.EmptyValues = nonNilStrings(result.EmptyValues)
	result.PotentialSecrets = nonNilStrings(result.PotentialSecrets)
	sort.Strings(result.EmptyValues)
	sort.Strings(result.PotentialSecrets)

	for _, key := range result.EmptyValues {
		report.Warnings = append(report.Warnings, fmt.Sprintf("Empty value: %s", key))
	}
	for _, key := range result.PotentialSecrets {
		report.Warnings = append(report.Warnings, fmt.Sprintf("Potential sensitive key: %s", key))
	}

	report.Status = reportStatus(0, len(report.Warnings))
	report.Summary["total_keys"] = result.TotalKeys
	report.Summary["empty_values"] = len(result.EmptyValues)
	report.Summary["potential_sensitive_keys"] = len(result.PotentialSecrets)
	report.Details = map[string]any{
		"empty_values":      result.EmptyValues,
		"potential_secrets": result.PotentialSecrets,
	}

	printJSONReport(report)
	return 0
}

func runDoctorJSON(envPath string, examplePath string) int {
	report := newJSONReport("doctor")
	report.TargetFile = envPath
	report.ExampleFile = examplePath

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
	securityWarnings := 0

	if !envExists {
		targetFileIssues = 1
		report.Errors = append(report.Errors, fmt.Sprintf("%s file missing", envPath))
	}

	if !exampleExists {
		exampleFileIssues = 1
		report.Warnings = append(report.Warnings, fmt.Sprintf("%s file missing", examplePath))
	}

	for _, key := range missingInEnv {
		report.Warnings = append(report.Warnings, fmt.Sprintf("Missing key in %s: %s", envPath, key))
	}

	if envTracked {
		securityWarnings = 1
		report.Warnings = append(report.Warnings, fmt.Sprintf("%s appears to be tracked by git", envPath))
	}

	exitCode := 0
	if len(missingInEnv) > 0 || !envExists {
		exitCode = 1
	}

	if exitCode != 0 {
		report.Status = "fail"
	} else {
		report.Status = reportStatus(0, len(report.Warnings))
	}

	report.Summary["target_file_issues"] = targetFileIssues
	report.Summary["example_file_issues"] = exampleFileIssues
	report.Summary["missing_keys"] = len(missingInEnv)
	report.Summary["security_warnings"] = securityWarnings
	report.Details = map[string]any{
		"target_file_exists":  envExists,
		"example_file_exists": exampleExists,
		"missing_in_env":      missingInEnv,
		"env_file_tracked":    envTracked,
	}

	printJSONReport(report)
	return exitCode
}

func runSecurityJSON(dirPath string, envPath string) int {
	report := newJSONReport("security")
	report.RootDirectory = dirPath
	report.TargetFile = envPath

	result, err := security.Run(dirPath, envPath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not run security checks: %v", err))
		report.Summary["secret_findings"] = 1
		report.Summary["warnings"] = 0
		printJSONReport(report)
		return 1
	}

	secretFindingCount := 0
	warningCount := 0

	if result.EnvFileTracked {
		report.Warnings = append(report.Warnings, fmt.Sprintf("%s appears to be tracked by git", envPath))
		warningCount++
	}

	for _, finding := range result.EnvFindings {
		report.Errors = append(report.Errors, fmt.Sprintf("Env secret: %s (%s)", finding.Location, finding.Kind))
		secretFindingCount++
	}
	for _, finding := range result.RepositoryFindings {
		report.Errors = append(report.Errors, fmt.Sprintf("Repository secret: %s (%s)", finding.Location, finding.Kind))
		secretFindingCount++
	}
	for _, finding := range result.HistoryFindings {
		report.Errors = append(report.Errors, fmt.Sprintf("Git history secret: %s (%s)", finding.Location, finding.Kind))
		secretFindingCount++
	}

	if !result.HistoryScanned {
		report.Warnings = append(report.Warnings, "Git history scan skipped")
		warningCount++
	}

	report.Status = reportStatus(secretFindingCount, warningCount)
	report.Summary["secret_findings"] = secretFindingCount
	report.Summary["warnings"] = warningCount
	report.Details = map[string]any{
		"env_findings":        result.EnvFindings,
		"repository_findings": result.RepositoryFindings,
		"history_findings":    result.HistoryFindings,
		"env_file_tracked":    result.EnvFileTracked,
		"history_scanned":     result.HistoryScanned,
	}

	printJSONReport(report)

	if secretFindingCount > 0 {
		return 1
	}
	return 0
}

func runLogScanJSON(dirPath string) int {
	report := newJSONReport("log-scan")
	report.RootDirectory = dirPath

	result, err := logscan.Run(dirPath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not scan logs in %s", dirPath))
		report.Summary["log_exposure_findings"] = 1
		printJSONReport(report)
		return 1
	}

	sort.Slice(result.Findings, func(i, j int) bool {
		if result.Findings[i].Location == result.Findings[j].Location {
			return result.Findings[i].Kind < result.Findings[j].Kind
		}
		return result.Findings[i].Location < result.Findings[j].Location
	})

	for _, finding := range result.Findings {
		report.Errors = append(report.Errors, fmt.Sprintf("Log exposure: %s (%s)", finding.Location, finding.Kind))
	}

	report.Status = reportStatus(len(result.Findings), 0)
	report.Summary["scanned_files"] = result.ScannedFilesCount
	report.Summary["log_exposure_findings"] = len(result.Findings)
	report.Details = map[string]any{
		"findings": result.Findings,
	}

	printJSONReport(report)

	if len(result.Findings) > 0 {
		return 1
	}
	return 0
}

func runDockerJSON(dockerfilePath string, envPath string) int {
	report := newJSONReport("docker")
	report.Dockerfile = dockerfilePath
	report.TargetFile = envPath

	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read %s", envPath))
		report.Summary["missing_docker_keys"] = 1
		printJSONReport(report)
		return 1
	}

	result, err := runtimecheck.ValidateDockerfile(dockerfilePath, envFile.Values)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read %s", dockerfilePath))
		report.Summary["missing_docker_keys"] = 1
		printJSONReport(report)
		return 1
	}

	for _, key := range result.MissingKeys {
		report.Errors = append(report.Errors, fmt.Sprintf("Missing Docker env key: %s", key))
	}

	report.Status = reportStatus(len(result.MissingKeys), 0)
	report.Summary["referenced_keys"] = len(result.ReferencedKeys)
	report.Summary["missing_docker_keys"] = len(result.MissingKeys)
	report.Details = map[string]any{
		"referenced_keys": result.ReferencedKeys,
		"missing_keys":    result.MissingKeys,
	}

	printJSONReport(report)

	if len(result.MissingKeys) > 0 {
		return 1
	}
	return 0
}

func runCIJSON(envPath string, examplePath string) int {
	report := newJSONReport("ci")
	report.TargetFile = envPath
	report.ExampleFile = examplePath

	lintResult, err := linter.Run(envPath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not lint %s", envPath))
		report.Summary["ci_errors"] = 1
		printJSONReport(report)
		return 1
	}

	lintResult.InvalidLines = nonNilStrings(lintResult.InvalidLines)
	sort.Strings(lintResult.InvalidLines)
	for _, issue := range lintResult.InvalidLines {
		report.Errors = append(report.Errors, fmt.Sprintf("Lint: %s", issue))
	}

	schema, err := loadOptionalTypeSchema()
	if err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("could not read type schema file: %v", err))
	}

	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read %s", envPath))
		report.Summary["ci_errors"] = len(report.Errors)
		printJSONReport(report)
		return 1
	}

	exampleFile, err := parser.ParseEnvFile(examplePath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read %s", examplePath))
		report.Summary["ci_errors"] = len(report.Errors)
		printJSONReport(report)
		return 1
	}

	result := validator.ValidateEnv(envFile, exampleFile, schema)

	sort.Strings(result.MissingKeys)
	sort.Strings(result.DuplicateKeys)
	sort.Strings(result.InvalidTypeValues)

	for _, key := range result.MissingKeys {
		report.Errors = append(report.Errors, fmt.Sprintf("Missing key: %s", key))
	}
	for _, key := range result.DuplicateKeys {
		report.Errors = append(report.Errors, fmt.Sprintf("Duplicate key: %s", key))
	}
	for _, issue := range result.InvalidTypeValues {
		report.Errors = append(report.Errors, fmt.Sprintf("Invalid type: %s", issue))
	}

	report.Status = reportStatus(len(report.Errors), 0)
	report.Summary["ci_errors"] = len(report.Errors)
	report.Details = map[string]any{
		"lint_issues":         lintResult.InvalidLines,
		"missing_keys":        result.MissingKeys,
		"duplicate_keys":      result.DuplicateKeys,
		"invalid_type_values": result.InvalidTypeValues,
		"type_schema_loaded":  len(schema) > 0,
	}

	printJSONReport(report)

	if len(report.Errors) > 0 {
		return 1
	}
	return 0
}

func runScanCodeJSON(dirPath string, envPath string) int {
	report := newJSONReport("scan-code")
	report.RootDirectory = dirPath
	report.TargetFile = envPath

	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not read %s", envPath))
		report.Summary["missing_env_keys"] = 1
		printJSONReport(report)
		return 1
	}

	result, err := codebase.Run(dirPath, envFile.Values)
	if err != nil {
		report.Status = "fail"
		report.Errors = append(report.Errors, fmt.Sprintf("could not scan codebase in %s", dirPath))
		report.Summary["missing_env_keys"] = 1
		printJSONReport(report)
		return 1
	}

	for _, mismatch := range result.NamingMismatches {
		report.Warnings = append(report.Warnings, fmt.Sprintf("Code uses %s but env file contains %s", mismatch.CodeKey, mismatch.EnvKey))
	}
	for _, key := range result.MissingInEnv {
		report.Errors = append(report.Errors, fmt.Sprintf("Missing env key for code usage: %s", key))
	}
	for _, key := range result.UnusedInEnv {
		report.Warnings = append(report.Warnings, fmt.Sprintf("Unused env key in codebase: %s", key))
	}

	report.Status = reportStatus(len(report.Errors), len(report.Warnings))
	report.Summary["scanned_files"] = result.ScannedFilesCount
	report.Summary["detected_env_usages"] = len(result.UsedKeys)
	report.Summary["naming_mismatches"] = len(result.NamingMismatches)
	report.Summary["missing_env_keys"] = len(result.MissingInEnv)
	report.Summary["unused_env_keys"] = len(result.UnusedInEnv)
	report.Details = map[string]any{
		"used_keys":         result.UsedKeys,
		"missing_in_env":    result.MissingInEnv,
		"unused_in_env":     result.UnusedInEnv,
		"naming_mismatches": result.NamingMismatches,
	}

	printJSONReport(report)

	if len(result.MissingInEnv) > 0 {
		return 1
	}
	return 0
}

func runEncrypt(inputPath string, outputPath string) int {
	err := encryption.EncryptFile(inputPath, outputPath, os.Getenv("ENVGUARD_KEY"))
	if err != nil {
		printfStatusLine("Error: could not encrypt %s: %v\n", inputPath, err)
		return 1
	}

	fmt.Println("Env Encryption Report")
	fmt.Println("---------------------")
	fmt.Printf("Input file: %s\n", inputPath)
	fmt.Printf("Output file: %s\n\n", outputPath)
	printStatusLine("[PASS] Env file encrypted")

	return 0
}

func runDecrypt(inputPath string, outputPath string) int {
	err := encryption.DecryptFile(inputPath, outputPath, os.Getenv("ENVGUARD_KEY"))
	if err != nil {
		printfStatusLine("Error: could not decrypt %s: %v\n", inputPath, err)
		return 1
	}

	fmt.Println("Env Decryption Report")
	fmt.Println("---------------------")
	fmt.Printf("Input file: %s\n", inputPath)
	fmt.Printf("Output file: %s\n\n", outputPath)
	printStatusLine("[PASS] Env file decrypted")

	return 0
}

func runDocker(dockerfilePath string, envPath string) int {
	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		printfStatusLine("Error: could not read %s\n", envPath)
		return 1
	}

	result, err := runtimecheck.ValidateDockerfile(dockerfilePath, envFile.Values)
	if err != nil {
		printfStatusLine("Error: could not read %s\n", dockerfilePath)
		return 1
	}

	fmt.Println("Docker Env Validation Report")
	fmt.Println("----------------------------")
	fmt.Printf("Dockerfile: %s\n", dockerfilePath)
	fmt.Printf("Env file: %s\n", envPath)
	fmt.Printf("Referenced keys: %d\n\n", len(result.ReferencedKeys))

	if len(result.MissingKeys) == 0 {
		printStatusLine("[PASS] Docker environment references are satisfied")
		fmt.Println("")
		fmt.Printf("Summary: %d missing Docker key(s)\n", len(result.MissingKeys))
		return 0
	}

	for _, key := range result.MissingKeys {
		printfStatusLine("[ERROR] Missing Docker env key: %s\n", key)
	}

	fmt.Println("")
	fmt.Printf("Summary: %d missing Docker key(s)\n", len(result.MissingKeys))

	return 1
}

func runCI(envPath string, examplePath string) int {
	fmt.Println("Env CI Report")
	fmt.Println("-------------")
	fmt.Printf("Target file: %s\n", envPath)
	fmt.Printf("Example file: %s\n\n", examplePath)

	errorCount := 0

	lintResult, err := linter.Run(envPath)
	if err != nil {
		printfStatusLine("[ERROR] could not lint %s\n", envPath)
		return 1
	}

	for _, issue := range lintResult.InvalidLines {
		printfStatusLine("[ERROR] Lint: %s\n", issue)
		errorCount++
	}

	schema := map[string]string{}
	loadedSchema, err := parser.LoadTypeSchema("examples/.env.types")
	if err != nil {
		if !os.IsNotExist(err) {
			printfStatusLine("[ERROR] could not read type schema file: %v\n", err)
			errorCount++
		}
	} else {
		schema = loadedSchema
	}

	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		printfStatusLine("[ERROR] could not read %s\n", envPath)
		return 1
	}

	exampleFile, err := parser.ParseEnvFile(examplePath)
	if err != nil {
		printfStatusLine("[ERROR] could not read %s\n", examplePath)
		return 1
	}

	result := validator.ValidateEnv(envFile, exampleFile, schema)

	sort.Strings(result.MissingKeys)
	sort.Strings(result.DuplicateKeys)
	sort.Strings(result.InvalidTypeValues)

	for _, key := range result.MissingKeys {
		printfStatusLine("[ERROR] Missing key: %s\n", key)
		errorCount++
	}

	for _, key := range result.DuplicateKeys {
		printfStatusLine("[ERROR] Duplicate key: %s\n", key)
		errorCount++
	}

	for _, issue := range result.InvalidTypeValues {
		printfStatusLine("[ERROR] Invalid type: %s\n", issue)
		errorCount++
	}

	if errorCount == 0 {
		printStatusLine("[PASS] CI environment validation passed")
		fmt.Println("")
		fmt.Printf("Summary: %d CI error(s)\n", errorCount)
		return 0
	}

	fmt.Println("")
	fmt.Printf("Summary: %d CI error(s)\n", errorCount)

	return 1
}

func runPreStart(envPath string, examplePath string, commandArgs []string) int {
	validationExitCode := runValidate(envPath, examplePath)
	if validationExitCode != 0 {
		fmt.Println("")
		printStatusLine("[ERROR] Pre-start validation failed")
		return validationExitCode
	}

	fmt.Println("")
	printfStatusLine("[RUN] %s\n", strings.Join(commandArgs, " "))

	command := exec.Command(commandArgs[0], commandArgs[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}

		printfStatusLine("Error: could not run command: %v\n", err)
		return 1
	}

	return 0
}

func runScanCode(dirPath string, envPath string) int {
	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		printfStatusLine("Error: could not read %s\n", envPath)
		return 1
	}

	result, err := codebase.Run(dirPath, envFile.Values)
	if err != nil {
		printfStatusLine("Error: could not scan codebase in %s\n", dirPath)
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
			printfStatusLine("[WARN] Code uses %s but env file contains %s\n", mismatch.CodeKey, mismatch.EnvKey)
		}
	}

	if len(result.MissingInEnv) > 0 {
		fmt.Println("\nUsed in code but missing in env file:")
		for _, key := range result.MissingInEnv {
			printfStatusLine("[ERROR] Missing env key for code usage: %s\n", key)
		}
	}

	if len(result.UnusedInEnv) > 0 {
		fmt.Println("\nPresent in env file but unused in code:")
		for _, key := range result.UnusedInEnv {
			printfStatusLine("[WARN] Unused env key in codebase: %s\n", key)
		}
	}

	if len(result.NamingMismatches) == 0 && len(result.MissingInEnv) == 0 && len(result.UnusedInEnv) == 0 {
		printStatusLine("\n[PASS] Codebase analysis found no issues")
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
		printfStatusLine("Error: could not read %s\n", envPath)
		return 1
	}

	outputPath := ".env.example"

	file, err := os.Create(outputPath)
	if err != nil {
		printfStatusLine("Error: could not create %s\n", outputPath)
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
			printStatusLine("Error: failed to write to file")
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
		printfStatusLine("Error: could not read %s\n", envPath)
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
		printfStatusLine("Error: could not open %s\n", examplePath)
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
			printStatusLine("Error: failed to write to file")
			return 1
		}
	}

	fmt.Println("Env Example Sync Report")
	fmt.Println("------------------------")
	fmt.Printf("Source file: %s\n", envPath)
	fmt.Printf("Synced file: %s\n\n", examplePath)

	if len(newKeys) == 0 {
		printStatusLine("[PASS] .env.example already up to date")
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
		commandArgs, jsonOutput, err := extractJSONFlag(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		if hasAllFlag(commandArgs) {
			envPath, examplePath, err := getValidatePaths(commandArgs)
			if err != nil {
				printfStatusLine("Error: %s\n", err)
				os.Exit(1)
			}
			if envPath != ".env" || examplePath != ".env.example" {
				printStatusLine("Error: --all cannot be used with --file or --example")
				os.Exit(1)
			}
			if jsonOutput {
				os.Exit(runValidateAllJSON())
			}
			os.Exit(runValidateAll())
		}
		envPath, examplePath, err := getValidatePaths(commandArgs)
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		if jsonOutput {
			os.Exit(runValidateJSON(envPath, examplePath))
		}
		os.Exit(runValidate(envPath, examplePath))
	case "lint":
		if hasHelpFlag(args[1:]) {
			printLintHelp()
			os.Exit(0)
		}
		commandArgs, jsonOutput, err := extractJSONFlag(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		envPath, err := getLintFilePath(commandArgs)
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		if jsonOutput {
			os.Exit(runLintJSON(envPath))
		}
		os.Exit(runLint(envPath))
	case "analyze":
		if hasHelpFlag(args[1:]) {
			printAnalyzeHelp()
			os.Exit(0)
		}
		commandArgs, jsonOutput, err := extractJSONFlag(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		envPath, err := getAnalyzeFilePath(commandArgs)
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		if jsonOutput {
			os.Exit(runAnalyzeJSON(envPath))
		}
		os.Exit(runAnalyze(envPath))
	case "doctor":
		if hasHelpFlag(args[1:]) {
			printDoctorHelp()
			os.Exit(0)
		}
		commandArgs, jsonOutput, err := extractJSONFlag(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		envPath, examplePath, err := getDoctorPaths(commandArgs)
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		if jsonOutput {
			os.Exit(runDoctorJSON(envPath, examplePath))
		}
		os.Exit(runDoctor(envPath, examplePath))
	case "scan-code":
		if hasHelpFlag(args[1:]) {
			printScanCodeHelp()
			os.Exit(0)
		}
		commandArgs, jsonOutput, err := extractJSONFlag(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		dirPath, envPath, err := getScanCodeOptions(commandArgs)
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		if jsonOutput {
			os.Exit(runScanCodeJSON(dirPath, envPath))
		}
		os.Exit(runScanCode(dirPath, envPath))
	case "security":
		if hasHelpFlag(args[1:]) {
			printSecurityHelp()
			os.Exit(0)
		}
		commandArgs, jsonOutput, err := extractJSONFlag(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		dirPath, envPath, err := getScanCodeOptions(commandArgs)
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		if jsonOutput {
			os.Exit(runSecurityJSON(dirPath, envPath))
		}
		os.Exit(runSecurity(dirPath, envPath))
	case "log-scan":
		if hasHelpFlag(args[1:]) {
			printLogScanHelp()
			os.Exit(0)
		}
		commandArgs, jsonOutput, err := extractJSONFlag(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		dirPath, err := getLogScanDir(commandArgs)
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		if jsonOutput {
			os.Exit(runLogScanJSON(dirPath))
		}
		os.Exit(runLogScan(dirPath))
	case "encrypt":
		if hasHelpFlag(args[1:]) {
			printEncryptHelp()
			os.Exit(0)
		}
		inputPath, outputPath, err := getCryptoPaths(args[1:], ".env", ".env.enc")
		if err != nil {
			printfStatusLine("Error: %s\n", err)
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
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runDecrypt(inputPath, outputPath))
	case "docker":
		if hasHelpFlag(args[1:]) {
			printDockerHelp()
			os.Exit(0)
		}
		commandArgs, jsonOutput, err := extractJSONFlag(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		dockerfilePath, envPath, err := getDockerOptions(commandArgs)
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		if jsonOutput {
			os.Exit(runDockerJSON(dockerfilePath, envPath))
		}
		os.Exit(runDocker(dockerfilePath, envPath))
	case "ci":
		if hasHelpFlag(args[1:]) {
			printCIHelp()
			os.Exit(0)
		}
		commandArgs, jsonOutput, err := extractJSONFlag(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		envPath, examplePath, err := getValidatePaths(commandArgs)
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		if jsonOutput {
			os.Exit(runCIJSON(envPath, examplePath))
		}
		os.Exit(runCI(envPath, examplePath))
	case "run":
		if hasHelpFlag(args[1:]) {
			printRunHelp()
			os.Exit(0)
		}
		envPath, examplePath, commandArgs, err := getRunOptions(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runPreStart(envPath, examplePath, commandArgs))
	case "generate-example":
		if hasHelpFlag(args[1:]) {
			printGenerateExampleHelp()
			os.Exit(0)
		}
		envPath, err := getLintFilePath(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
			os.Exit(1)
		}
		os.Exit(runGenerateExample(envPath))
	case "sync-example":
		if hasHelpFlag(args[1:]) {
			printSyncExampleHelp()
			os.Exit(0)
		}
		envPath, err := getLintFilePath(args[1:])
		if err != nil {
			printfStatusLine("Error: %s\n", err)
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
