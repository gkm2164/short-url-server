package model

import (
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	ShortenId   string `json:"shorten_id"`
	Url         string `json:"url"`
	AccessCount int64  `json:"access_count"`
}
