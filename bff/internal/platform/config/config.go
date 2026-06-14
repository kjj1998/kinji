package config

import "os"

type Config struct {
	Port            string
	Env             string
	CORSOrigin      string
	DBDriver        string // "dynamo" | "sqlite"
	SQLitePath      string
	AnthropicApiKey string
	AnthropicModel  string
}

func Load() Config {
	return Config{
		AnthropicApiKey: getEnv("ANTHROPIC_API_KEY", ""),
		AnthropicModel:  getEnv("ANTHROPIC_MODEL", "claude-sonnet-4-6"),
		Port:            getEnv("PORT", "8080"),
		Env:             getEnv("ENV", "development"),
		CORSOrigin:      getEnv("CORS_ORIGIN", "http://localhost:5173"),
		DBDriver:        getEnv("DB_DRIVER", "sqlite"),
		SQLitePath:      getEnv("SQLITE_PATH", "kinji.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
