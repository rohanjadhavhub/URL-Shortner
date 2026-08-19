package shortner

import (
		"math/rand"
		"errors"
		"log"
		"fmt"

	lru "github.com/hashicorp/golang-lru/v2"
		"github.com/jackc/pgx/v5/pgconn"
		"gorm.io/gorm"

		"url-shortener/models"

)

type Cache struct {
	lru *lru.Cache[string, string]
}

const defaultCacheSize = 1000
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
		defaultCodeLength = 7
		maxRetries = 5
)

func NewCache() *Cache {
	c, err := lru.New[string, string](defaultCacheSize)
	if err != nil {
		panic(err)
	}

	return &Cache{lru: c}
}

func (c *Cache) Get(code string) (string, bool) {
	return c.lru.Get(code)
}

func (c *Cache) Set(code, longURL string) {
	c.lru.Add(code, longURL)
}

func GenerateCode(length int) string{
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func CreateShortURL(db *gorm.DB, longURL string) (*models.URL, error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		url := &models.URL{
				LongURL:   longURL,
				ShortURL: GenerateCode(defaultCodeLength),
		}

		err := db.Create(url).Error
		if err == nil {
			return url, nil
		}

		if isUniqueViolation(err) {
			continue
		}

		return nil, err
	}
	return nil, fmt.Errorf("failed to generate unique short code after %d attempts", maxRetries)
}


func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}


func GetByShortCodeCached(db *gorm.DB, cache *Cache, code string) (*models.URL, error) {
	if longURL, ok := cache.Get(code); ok {
		log.Printf("cache HIT for code=%s", code)
		return &models.URL{ShortURL: code, LongURL: longURL}, nil
	}

	log.Printf("cache MISS for code=%s", code)

	url, err := GetByShortCode(db, code)
	if err != nil {
		return nil, err
	}

	cache.Set(url.ShortURL, url.LongURL)
	return url, nil
}
func GetByShortCode(db *gorm.DB, code string) (*models.URL, error) {
	var url models.URL
	if err := db.Take(&url, "short_url = ?", code).Error; err != nil {
		return nil, err
	}
	return &url, nil
}