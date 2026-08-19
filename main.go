package main

import(
	"fmt"
	"log"

	"url-shortener/config"
	"url-shortener/db"
	"url-shortener/models"
	"url-shortener/shortner"

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

	if err := database.AutoMigrate(&models.URL{}); err != nil {
    	log.Fatalf("failed to migrate database: %v", err)
	}

	url, err := shortner.CreateShortURL(database, "https://github.com/Anshgrover23/vouch")
	if err != nil {
		log.Fatalf("failed to create short url: %v", err)
	}
	fmt.Printf("Created: %s -> %s (id=%d)\n", url.ShortURL, url.LongURL, url.ID)

	found, err := shortner.GetByShortCode(database, url.ShortURL)
	if err != nil {
		log.Fatalf("lookup failed: %v", err)
	}
	fmt.Printf("Found: %s -> %s\n", found.ShortURL, found.LongURL)

	fmt.Printf("Starting in %s mode on port %s\n", cfg.Env, cfg.Port)
}