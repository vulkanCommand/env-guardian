package doctor

import (
	"os"

	"github.com/vulkanCommand/env-guardian/internal/parser"
	"github.com/vulkanCommand/env-guardian/internal/security"
	"github.com/vulkanCommand/env-guardian/internal/validator"
)

type DoctorResult struct {
	EnvFileExists     bool
	ExampleFileExists bool
	EnvFileTracked    bool
	MissingInEnv      []string
}

func Run() DoctorResult {
	result := DoctorResult{}

	if _, err := os.Stat(".env"); err == nil {
		result.EnvFileExists = true
	}

	if _, err := os.Stat(".env.example"); err == nil {
		result.ExampleFileExists = true
	}

	result.EnvFileTracked = security.IsGitTracked(".", ".env")

	if result.EnvFileExists && result.ExampleFileExists {
		envFile, err1 := parser.ParseEnvFile(".env")
		exampleFile, err2 := parser.ParseEnvFile(".env.example")

		if err1 == nil && err2 == nil {
			validationResult := validator.ValidateEnv(envFile, exampleFile, map[string]string{})
			result.MissingInEnv = validationResult.MissingKeys
		}
	}

	return result
}
