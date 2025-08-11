package main

import (
	"log"
	"net/http"

	"github.com/go-squad-5/pdf-generator/internal/handler"
	"github.com/go-squad-5/pdf-generator/internal/repository"
	"github.com/go-squad-5/pdf-generator/internal/router"
	"github.com/go-squad-5/pdf-generator/internal/service"
)

func main() {
	db, err := repository.InitDB("./database.sqlite")
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer db.Close()
	log.Println("Database initialized successfully.")

	userRepo := repository.NewUserRepository(db)
	markRepo := repository.NewMarkRepository(db)
	pdfService := service.NewPDFService(userRepo, markRepo)
	pdfHandler := handler.NewPDFHandler(pdfService)

	r := router.NewRouter(pdfHandler)

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
