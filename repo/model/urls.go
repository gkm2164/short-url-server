package model

import "time"

type Url struct {
	Id        string    `json:"id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}
