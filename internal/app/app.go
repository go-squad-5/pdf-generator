package app

import "github.com/go-squad-5/pdf-generator/internal/config"

type App struct {
	config *config.Config
	// and more fields... db, logger, etc.
}

func NewApp(config *config.Config) *App {
	return &App{
		config: config,
		// Initialize other fields like db, logger, etc.
	}
}
