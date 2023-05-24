package server

import (
	"net/http"
)

type server struct {
	addr    string
	handler http.Handler
}

// New конструктор сервера
func New(addr string, handler http.Handler) *server {
	return &server{
		addr:    addr,
		handler: handler,
	}
}

// Start запуск сервера
func (s *server) Start() error {
	err := http.ListenAndServe(s.addr, s.handler)

	return err
}
