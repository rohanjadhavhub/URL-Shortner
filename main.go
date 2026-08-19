package main

import(
	"fmt"
	"log"

	"url-shortener/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config:", err)
	}
	
	fmt.Printf("Starting in %s mode on port %s\n", cfg.Env, cfg.Port)
}