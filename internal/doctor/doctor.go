package doctor

import "os"

type DoctorResult struct {
	EnvFileExists     bool
	ExampleFileExists bool
}

func Run() DoctorResult {
	result := DoctorResult{}

	if _, err := os.Stat(".env"); err == nil {
		result.EnvFileExists = true
	}

	if _, err := os.Stat(".env.example"); err == nil {
		result.ExampleFileExists = true
	}

	return result
}