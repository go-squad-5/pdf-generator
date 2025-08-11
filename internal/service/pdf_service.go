package service

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/go-squad-5/pdf-generator/internal/models"
	"github.com/go-squad-5/pdf-generator/internal/repository"
	"github.com/jung-kurt/gofpdf"
)

type PDFService struct {
	sessionRepo *repository.SessionRepository
	quizzesRepo *repository.QuizzesRepository
}

func NewPDFService(sessionRepo *repository.SessionRepository, quizzesRepo *repository.QuizzesRepository) *PDFService {
	return &PDFService{sessionRepo: sessionRepo, quizzesRepo: quizzesRepo}
}

func (s *PDFService) GenerateQuizReport(sessionID string) ([]byte, error) {
	var wg sync.WaitGroup
	var session *models.Session
	var quizzes []models.Quiz
	var sessionErr, quizzesErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		session, sessionErr = s.sessionRepo.GetSessionByID(sessionID)
	}()
	go func() {
		defer wg.Done()
		quizzes, quizzesErr = s.quizzesRepo.GetQuizzesBySessionID(sessionID)
	}()
	wg.Wait()

	if sessionErr != nil || quizzesErr != nil {
		return nil, fmt.Errorf("failed to fetch data: sessionErr=%v, quizzesErr=%v", sessionErr, quizzesErr)
	}
	if session == nil {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}

	return s.createPDF(session, quizzes)
}

func (s *PDFService) createPDF(session *models.Session, quizzes []models.Quiz) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "Quiz Performance Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "User Email:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, session.Email)
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Final Score:")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 100, 0)
	pdf.Cell(0, 8, fmt.Sprintf("%d / %d", session.Score, len(quizzes)))
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(15)

	for i, quiz := range quizzes {
		pdf.SetFont("Arial", "B", 10)
		pdf.MultiCell(0, 5, fmt.Sprintf("%d. %s", i+1, quiz.QuestionData.Question), "", "L", false)
		pdf.Ln(2)

		pdf.SetFont("Arial", "", 9)
		if quiz.IsCorrect {
			pdf.SetFillColor(200, 255, 200)
			pdf.Cell(40, 5, "Your Answer (Correct):")
		} else {
			pdf.SetFillColor(255, 200, 200)
			pdf.Cell(40, 5, "Your Answer (Incorrect):")
		}
		pdf.CellFormat(0, 5, quiz.Answer, "", 0, "L", true, 0, "")
		pdf.Ln(5)

		if !quiz.IsCorrect {
			pdf.SetFillColor(230, 230, 230)
			pdf.Cell(40, 5, "Correct Answer:")
			pdf.CellFormat(0, 5, quiz.QuestionData.Answer, "", 0, "L", true, 0, "")
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
