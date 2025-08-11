package service

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/go-squad-5/pdf-generator/internal/models"
	"github.com/go-squad-5/pdf-generator/internal/repository"
	"github.com/jung-kurt/gofpdf"
)

type PDFService struct {
	sessionRepo *repository.SessionRepository
	attemptRepo *repository.QuizAttemptRepository
}

func NewPDFService(sessionRepo *repository.SessionRepository, attemptRepo *repository.QuizAttemptRepository) *PDFService {
	return &PDFService{sessionRepo: sessionRepo, attemptRepo: attemptRepo}
}

func (s *PDFService) GenerateQuizReport(sessionID int) ([]byte, error) {
	var wg sync.WaitGroup
	var session *models.Session
	var attempts []models.QuizAttempt
	var sessionErr, attemptsErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		session, sessionErr = s.sessionRepo.GetSessionByID(sessionID)
	}()

	go func() {
		defer wg.Done()
		attempts, attemptsErr = s.attemptRepo.GetAttemptsBySessionID(sessionID)
	}()

	log.Println("Main thread is waiting for goroutines to finish...")
	wg.Wait()
	log.Println("All goroutines finished.")

	if sessionErr != nil {
		return nil, fmt.Errorf("failed to fetch session data: %w", sessionErr)
	}
	if attemptsErr != nil {
		return nil, fmt.Errorf("failed to fetch attempts data: %w", attemptsErr)
	}
	if session == nil {
		return nil, fmt.Errorf("session with ID %d not found", sessionID)
	}

	return s.createPDF(session, attempts)
}

func (s *PDFService) createPDF(session *models.Session, attempts []models.QuizAttempt) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "Quiz Performance Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Student Name:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName))
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Session ID:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("%d", session.ID))
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Final Score:")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 100, 0)
	pdf.Cell(0, 8, fmt.Sprintf("%d / %d", session.TotalMarks, len(attempts)))
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(15)

	for i, attempt := range attempts {
		pdf.SetFont("Arial", "B", 10)
		pdf.MultiCell(0, 5, fmt.Sprintf("%d. %s", i+1, attempt.Question.QuestionText), "", "L", false)
		pdf.Ln(2)

		pdf.SetFont("Arial", "", 9)
		isCorrect := attempt.ChosenOption == attempt.Question.CorrectOption
		if isCorrect {
			pdf.SetFillColor(200, 255, 200)
			pdf.Cell(35, 5, "Your Answer (Correct):")
		} else {
			pdf.SetFillColor(255, 200, 200)
			pdf.Cell(35, 5, "Your Answer (Incorrect):")
		}
		pdf.CellFormat(0, 5, fmt.Sprintf("'%s'", attempt.ChosenOption), "", 0, "L", true, 0, "")
		pdf.Ln(5)

		if !isCorrect {
			pdf.SetFillColor(230, 230, 230)
			pdf.Cell(35, 5, "Correct Answer:")
			pdf.CellFormat(0, 5, fmt.Sprintf("'%s'", attempt.Question.CorrectOption), "", 0, "L", true, 0, "")
			pdf.Ln(5)
		}
		pdf.Ln(5)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
