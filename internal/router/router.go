package router

import (
	"net/http"

	"github.com/go-squad-5/pdf-generator/internal/app"
	"github.com/go-squad-5/pdf-generator/internal/handlers"
)

type Router struct {
	app *app.App
}

func NewRouter(app *app.App) *Router {
	return &Router{
		app: app,
	}
}

func (r *Router) Route() *http.ServeMux {
	mux := http.NewServeMux()
	// Initialize the handler with the app instance
	h := handlers.NewHandler(r.app)

	// docs
	mux.HandleFunc("/docs/openapi.yaml", h.ServeOpenAPISpec)
	mux.Handle("/docs/", h.SwaggerUIHandler())

	// Versioned routes
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", r.v1Routes()))

	// add more versioned routes as needed

	return mux
}

func (r *Router) v1Routes() http.Handler {
	// Example of adding a versioned route
	mux := http.NewServeMux()

	h := handlers.NewHandler(r.app)

	// Health check route
	mux.HandleFunc("/health", h.HealthCheck)

	// Add more v1 routes here

	return mux
}
