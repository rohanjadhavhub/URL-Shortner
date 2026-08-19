package config

import(
	"fmt"
	"os"

	"github.com/joho/godotenv"
)
type Config struct {
	Env         string
	DatabaseURL string
	Port        string
}


func Load() (*Config, error){

	_ = godotenv.Load()

	cfg := &Config{
		Env:           getEnv("APP_ENV", "development"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		Port:          getEnv("PORT", "8080"),
	}

	if cfg.DatabaseURL=="" {
		return nil, fmt.Errorf("DatabaseURL needed But not set yet.")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != ""{
		return val
	}
	
	return fallback
}