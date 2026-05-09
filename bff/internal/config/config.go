package config

import "os"

type Config struct {
	Port           string
	DynamoEndpoint string
	DynamoRegion   string
	DynamoTable    string
	ParserGRPCAddr string
	Env            string
	CORSOrigin     string
}

func Load() Config {
	return Config{
		Port:           getEnv("PORT", "8080"),
		DynamoEndpoint: getEnv("DYNAMO_ENDPOINT", "http://localhost:8000"),
		DynamoRegion:   getEnv("DYNAMO_REGION", "ap-southeast-1"),
		DynamoTable:    getEnv("DYNAMO_TABLE", "kinji"),
		ParserGRPCAddr: getEnv("PARSER_GRPC_ADDR", "localhost:50051"),
		Env:            getEnv("ENV", "development"),
		CORSOrigin:     getEnv("CORS_ORIGIN", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
