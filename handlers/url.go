package handlers

import (
		"encoding/json"
		"errors"
		"net/http"
		"net/url"

		"github.com/go-chi/chi/v5"
		"gorm.io/gorm"

		"url-shortener/shortner"
)


type Handler struct {
	DB    *gorm.DB
	Cache *shortner.Cache
}

type ShortenRequest struct {
	URL string `json: "url"`
}

type ShortenResponse struct {
	ShortURL string `json: "short_url"`
	LongURL string `json: "long_url"`
}

func validateURL(raw string) error {
	if raw == "" {
		return errors.New("url is required")
	}

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return errors.New("url is not well-formed")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("url must use http or https")
	}

	if u.Host == "" {
		return errors.New("url must include a host")
	}
	return nil
}



func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return 
	}

	if err := validateURL(req.URL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	url, err := shortner.CreateShortURL(h.DB, req.URL)
	if err != nil {
		http.Error(w, "failed to create a short url", http.StatusInternalServerError)
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

	url, err := shortner.GetByShortCodeCached(h.DB, h.Cache, code)
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