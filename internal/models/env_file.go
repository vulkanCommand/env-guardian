package models

type EnvFile struct {
	Values     map[string]string
	Duplicates map[string]int
}