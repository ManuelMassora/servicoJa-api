package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseDSN        string
	ServerPort         string
}

func LoadConfig(envPath string) (*Config, error) {

	err := godotenv.Load(envPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("erro ao carregar arquivo .env: %w", err)
	}

	dbDSN := os.Getenv("DATABASE_DSN")
	if dbDSN == "" {
		log.Printf("DATABASE_DSN not set, using default: %s", dbDSN)
	}

	serverPort := os.Getenv("SERVER_PORT")

	if dbDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN environment variable is not set")
	}

	if dbDSN == "" {
		return nil, fmt.Errorf("configurações essenciais (DATABASE_DSN, S3_BUCKET_NAME, AWS_REGION) não estão definidas")
	}

	cfg := &Config{
		DatabaseDSN:        dbDSN,
		ServerPort:         serverPort,
	}
	return cfg, nil
}