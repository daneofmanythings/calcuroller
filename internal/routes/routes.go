package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func InitRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{}))

	return router
}
