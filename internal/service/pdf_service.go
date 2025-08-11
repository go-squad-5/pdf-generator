package service

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-squad-5/pdf-generator/internal/models"
	"github.com/go-squad-5/pdf-generator/internal/repository"
	"github.com/jung-kurt/gofpdf"
)

type PDFService struct {
	userRepo *repository.UserRepository
	markRepo *repository.MarkRepository
}

func NewPDFService(userRepo *repository.UserRepository, markRepo *repository.MarkRepository) *PDFService {
	return &PDFService{
		userRepo: userRepo,
		markRepo: markRepo,
	}
}

func (s *PDFService) GenerateUserReport(userID int) ([]byte, error) {
	var wg sync.WaitGroup
	var user *models.User
	var marks []models.Mark
	var userErr, marksErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Printf("Goroutine started for fetching user ID: %d", userID)
		user, userErr = s.userRepo.GetUserByID(userID)
	}()

	go func() {
		defer wg.Done()
		log.Printf("Goroutine started for fetching marks for user ID: %d", userID)
		marks, marksErr = s.markRepo.GetMarksByUserID(userID)
	}()

	log.Println("Main thread is waiting for goroutines to finish...")
	wg.Wait()
	log.Println("All goroutines finished.")

	if userErr != nil {
		log.Printf("Error fetching user data: %v", userErr)
		return nil, fmt.Errorf("failed to fetch user data: %w", userErr)
	}
	if marksErr != nil {
		log.Printf("Error fetching marks data: %v", marksErr)
		return nil, fmt.Errorf("failed to fetch marks data: %w", marksErr)
	}

	if user == nil {
		return nil, fmt.Errorf("user with ID %d not found", userID)
	}

	return s.createPDF(user, marks)
}

func (s *PDFService) createPDF(user *models.User, marks []models.Mark) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 24)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(0, 20, "Student Report Card", "1", 1, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Student Name:")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 10, fmt.Sprintf("%s %s", user.FirstName, user.LastName), "", 1, "", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Student ID:")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 10, fmt.Sprintf("%d", user.ID), "", 1, "", false, 0, "")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Marks Report")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(130, 10, "Subject", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 10, "Score", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 12)
	totalScore := 0
	for _, mark := range marks {
		pdf.CellFormat(130, 10, mark.Subject, "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 10, fmt.Sprintf("%d", mark.Score), "1", 1, "C", false, 0, "")
		totalScore += mark.Score
	}

	average := float64(totalScore) / float64(len(marks))
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(130, 10, "Average Score", "1", 0, "R", false, 0, "")
	pdf.CellFormat(60, 10, fmt.Sprintf("%.2f", average), "1", 1, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	footerText := fmt.Sprintf("Report generated on %s", time.Now().Format("2006-01-02 15:04:05"))
	pdf.CellFormat(0, 10, footerText, "", 0, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
