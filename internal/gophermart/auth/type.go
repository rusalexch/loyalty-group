package auth

import "github.com/go-chi/chi"

type Config struct {
	Mux         *chi.Mux
	UserService userService
	JwtSecret   string
}

type authModule struct {
	mux         *chi.Mux
	userService userService
	jwtSecret   string
}

type signup struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
