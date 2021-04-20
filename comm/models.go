package comm

import "time"

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type FailureResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AllUrlResponse struct {
	Status string      `json:"status"`
	Urls   []UrlEntity `json:"urls"`
}

type GetUrlResponse struct {
	Status string    `json:"status"`
	Url    UrlEntity `json:"url"`
}

type UrlEntity struct {
	Id        string    `json:"id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUrlRequest struct {
	Url string `json:"url"`
}

type CreateUrlResponse struct {
	Status string `json:"status"`
	Id     string `json:"id"`
}

type UpdateUrlRequest struct {
	Url string `json:"url"`
}

type UpdateUrlResponse struct {
	Status  string `json:"status"`
	Updated bool   `json:"url"`
}
