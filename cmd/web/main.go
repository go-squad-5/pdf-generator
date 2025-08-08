package main

import (
	"log"
	"net/http"

	"github.com/go-squad-5/pdf-generator/internal/app"
	"github.com/go-squad-5/pdf-generator/internal/config"
	"github.com/go-squad-5/pdf-generator/internal/router"
)

func main() {
	// Load configuration
	cfg := config.Load()
	// Initialize the application
	app := app.NewApp(cfg)
	// Initialize the router with the app instance
	route := router.NewRouter(app).Route()
	// Start the HTTP server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", route))
}
