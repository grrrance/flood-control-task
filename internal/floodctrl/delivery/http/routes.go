package http

import (
	"github.com/go-chi/chi/v5"
	"task/internal/floodctrl"
)

func MapFloodRoutes(handlers floodctrl.Handlers) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/user/{id}", handlers.TriggerUser())
	}
}
