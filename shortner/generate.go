package shortner

import (
	"math/rand"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"url-shortener/models"

)
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
		defaultCodeLength = 7
		maxRetries = 5
)

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

func GetByShortCode(db *gorm.DB, code string) (*models.URL, error) {
	var url models.URL
	if err := db.Take(&url, "short_url = ?", code).Error; err != nil {
		return nil, err
	}
	return &url, nil
}