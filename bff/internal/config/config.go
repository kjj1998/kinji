package config

import "os"

type Config struct {
	Port           string
	ParserGRPCAddr string
	Env            string
	CORSOrigin     string
	DBDriver       string // "dynamo" | "sqlite"
	SQLitePath     string
}

func Load() Config {
	return Config{
		Port:           getEnv("PORT", "8080"),
		ParserGRPCAddr: getEnv("PARSER_GRPC_ADDR", "localhost:50051"),
		Env:            getEnv("ENV", "development"),
		CORSOrigin:     getEnv("CORS_ORIGIN", "http://localhost:5173"),
		DBDriver:       getEnv("DB_DRIVER", "sqlite"),
		SQLitePath:     getEnv("SQLITE_PATH", "kinji.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
