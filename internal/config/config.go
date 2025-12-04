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
	ServerHost         string
	SupabaseURL        string
	SupabaseKey        string
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
	serverHost := os.Getenv("SERVER_HOST")
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")

	if dbDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN environment variable is not set")
	}

	if supabaseURL == "" || supabaseKey == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_SERVICE_KEY environment variables must be set")
	}

	cfg := &Config{
		DatabaseDSN:        dbDSN,
		ServerPort:         serverPort,
		ServerHost:         serverHost,
		SupabaseURL:        supabaseURL,
		SupabaseKey:        supabaseKey,
	}
	return cfg, nil
}