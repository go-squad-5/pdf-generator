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
	// Step 1: Connect to your MySQL database.
	db, err := repository.InitDB()
	if err != nil {
		log.Fatalf("Could not initialize and connect to MySQL: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to MySQL database.")

	// Step 2: Initialize repositories
	sessionRepo := repository.NewSessionRepository(db)
	quizzesRepo := repository.NewQuizzesRepository(db)

	// Step 3: Initialize services
	pdfService := service.NewPDFService(sessionRepo, quizzesRepo)
	emailService := service.NewEmailService(sessionRepo, quizzesRepo)

	// Step 4: Initialize handlers
	pdfHandler := handler.NewPDFHandler(pdfService)
	emailHandler := handler.NewEmailHandler(emailService)

	// Step 5: Setup Router
	r := router.NewRouter(pdfHandler, emailHandler)

	// Step 6: Start Server
	port := ":8070"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
