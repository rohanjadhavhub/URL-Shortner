package main

import(
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"url-shortener/config"
	"url-shortener/db"
	"url-shortener/models"
	"url-shortener/handlers"
	"url-shortener/shortner"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config:%v", err)
	}

	database, err := db.ConnectDB(cfg.DatabaseURL, cfg.Env)
	if err != nil {
		log.Fatalf("Failed to connect DB:%v", err)
	}

	if err := database.AutoMigrate(&models.URL{}); err != nil {
    	log.Fatalf("failed to migrate database: %v", err)
	}

	h := &handlers.Handler{
			DB: database,
			Cache: shortner.NewCache(),
		}

	r := chi.NewRouter()
	r.Post("/shorten", h.Shorten)
	r.Get("/{code}", h.Redirect)

	log.Printf("starting server in %s mode on port %s\n", cfg.Env, cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}