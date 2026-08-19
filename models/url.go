package models

import "time"

type URL struct {
	ID        uint     `gorm: "primaryKey"`
	LongURL   string   `gorm: "type:text; not null"`
	ShortURL  string   `gorm: "uniqueIndex; size: 10; not null"`
	CreatedAt time.Time
}

