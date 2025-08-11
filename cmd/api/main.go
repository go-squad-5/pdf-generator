package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-squad-5/pdf-generator/internal/handler"
	"github.com/go-squad-5/pdf-generator/internal/repository"
	"github.com/go-squad-5/pdf-generator/internal/router"
	"github.com/go-squad-5/pdf-generator/internal/service"
)

func main() {
	os.Remove("./database.sqlite")

	db, err := repository.InitDB("./database.sqlite")
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer db.Close()
	log.Println("Database initialized successfully.")

	sessionRepo := repository.NewSessionRepository(db)
	attemptRepo := repository.NewQuizAttemptRepository(db)

	pdfService := service.NewPDFService(sessionRepo, attemptRepo)
	emailService := service.NewEmailService(sessionRepo, attemptRepo)

	pdfHandler := handler.NewPDFHandler(pdfService)
	emailHandler := handler.NewEmailHandler(emailService)

	r := router.NewRouter(pdfHandler, emailHandler)

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
