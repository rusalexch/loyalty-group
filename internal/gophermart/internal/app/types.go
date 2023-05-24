package app

import (
	"net/http"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"-"`
}

type CreateUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	ID         string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type Handler struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

type Middleware func(h http.Handler) http.Handler
