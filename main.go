package main

import(
	"fmt"
	"log"

	"url-shortener/config"
	"url-shortener/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config:", err)
	}

	database, err := db.ConnectDB(cfg.DatabaseURL, cfg.Env)
	if err != nil {
		log.Fatalf("Failed to connect DB:", err)
	}

	_ = database

	
	fmt.Printf("Starting in %s mode on port %s\n", cfg.Env, cfg.Port)
}