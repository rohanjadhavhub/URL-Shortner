package handlers

import (
		"encoding/json"
		"errors"
		"net/http"

		"github.com/go-chi/chi/v5"
		"gorm.io/gorm"

		"url-shortener/shortner"
)


type Handler struct {
	DB *gorm.DB
}

type ShortenRequest struct {
	URL string `json: "url"`
}

type ShortenResponse struct {
	ShortURL string `json: "short_url"`
	LongURL string `json: "long_url"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return 
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	
	url, err := shortner.CreateShortURL(h.DB, req.URL)
	if err != nil {
		http.Error(w, "failed to create a short url", http.StatusBadRequest)
		return
	}

	resp := ShortenResponse{
			ShortURL: url.ShortURL,
			LongURL: url.LongURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, err := shortner.GetByShortCode(h.DB, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			http.Error(w, "short url not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to lookup a short url", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url.LongURL, http.StatusFound)

}