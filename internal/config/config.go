package config

import "os"

// Config holds application configuration values.
type Config struct {
	GitHubToken string
	GitHubOrg   string
}

// LoadFromEnv populates Config from environment variables.
func LoadFromEnv() Config {
	return Config{
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
		GitHubOrg:   os.Getenv("GITHUB_ORG"),
	}
}
