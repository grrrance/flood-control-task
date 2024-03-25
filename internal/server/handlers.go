package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	floodHttp "task/internal/floodctrl/delivery/http"
	"task/internal/floodctrl/repository"
	"task/internal/floodctrl/usecase"
)

func (s *Server) MapHandlers(r *chi.Mux) error {

	floodRepo := repository.NewFloodRepository(s.db)
	floodUC := usecase.NewFloodUC(floodRepo, s.cfg)
	floodHandlers := floodHttp.NewFloodHandlers(floodUC, s.logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(corsHandle())

	r.Route("/flood", floodHttp.MapFloodRoutes(floodHandlers))

	return nil
}

func corsHandle() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
