package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
}

func startMainRouter() *Server {
	s := &Server{
		router: chi.NewRouter(),
	}
	s.setupMiddleware()
	s.setupRoutes()
	return s
}

func (s *Server) setupMiddleware() {
	// Basic middleware
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
}

func (s *Server) setupRoutes() {
	s.router.Get("/health", s.handleHealth())
	fs := http.FileServer(http.Dir(".well-known"))
	s.router.Handle("/.well-known/*", http.StripPrefix("/.well-known/", fs))
}

func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func (s *Server) Start(ctx context.Context) error {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return http.ListenAndServe(":"+port, s.router)
}
