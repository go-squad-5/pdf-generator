package service

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"sync"

	"github.com/go-squad-5/pdf-generator/internal/models"
	"github.com/go-squad-5/pdf-generator/internal/repository"
	"gopkg.in/gomail.v2"
)

const (
	smtpHost = "localhost"
	smtpPort = 1025
)

type EmailService struct {
	sessionRepo *repository.SessionRepository
	quizzesRepo *repository.QuizzesRepository
	dialer      *gomail.Dialer
}

func NewEmailService(sessionRepo *repository.SessionRepository, quizzesRepo *repository.QuizzesRepository) *EmailService {
	d := gomail.NewDialer(smtpHost, smtpPort, "", "")
	return &EmailService{sessionRepo: sessionRepo, quizzesRepo: quizzesRepo, dialer: d}
}

func (s *EmailService) SendQuizReportByEmail(sessionID string) error {
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
		return fmt.Errorf("failed to fetch data: sessionErr=%v, quizzesErr=%v", sessionErr, quizzesErr)
	}
	if session == nil {
		return fmt.Errorf("session with ID %s not found", sessionID)
	}

	const questionsPerEmail = 10
	var emailWg sync.WaitGroup

	for i := 0; i < len(quizzes); i += questionsPerEmail {
		end := i + questionsPerEmail
		if end > len(quizzes) {
			end = len(quizzes)
		}
		paginatedQuizzes := quizzes[i:end]
		pageNumber := (i / questionsPerEmail) + 1
		totalPages := (len(quizzes) + questionsPerEmail - 1) / questionsPerEmail

		emailWg.Add(1)
		go func(pQuizzes []models.Quiz, pNum, tPages int) {
			defer emailWg.Done()
			s.sendSingleEmailPart(session, pQuizzes, pNum, tPages)
		}(paginatedQuizzes, pageNumber, totalPages)
	}

	emailWg.Wait()
	return nil
}

func (s *EmailService) sendSingleEmailPart(session *models.Session, quizzesChunk []models.Quiz, pageNum, totalPages int) {
	body, err := s.parseEmailTemplate(session, quizzesChunk, pageNum, totalPages)
	if err != nil {
		log.Printf("ERROR: Could not parse email template for session %s: %v", session.SessionID, err)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "quiz-system@university.com")
	m.SetHeader("To", session.Email)
	m.SetHeader("Subject", fmt.Sprintf("Detailed Quiz Report for Session #%s (Part %d/%d)", session.SessionID, pageNum, totalPages))
	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("ERROR: Failed to send email for session %s part %d: %v", session.SessionID, pageNum, err)
	}
}

func (s *EmailService) parseEmailTemplate(session *models.Session, quizzes []models.Quiz, pageNum, totalPages int) (string, error) {
	const templateStr = `
    <!DOCTYPE html>
    <html>
    <body>
        <h1>Detailed Quiz Report</h1>
        <h2>Session #{{.Session.SessionID}} - Part {{.PageNum}} of {{.TotalPages}}</h2>
        <p>Dear user, here are the details for this part of your quiz:</p>
        {{range .Quizzes}}
        <div>
            <p><b>Question:</b> {{.QuestionData.Question}}</p>
            <p>Your Answer: <b>{{.Answer}}</b></p>
            {{if not .IsCorrect}}
            <p>Correct Answer: <b>{{.QuestionData.Answer}}</b></p>
            {{end}}
        </div>
        <hr>
        {{end}}
    </body>
    </html>`

	t, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", err
	}

	data := struct {
		Session    *models.Session
		Quizzes    []models.Quiz
		PageNum    int
		TotalPages int
	}{
		Session:    session,
		Quizzes:    quizzes,
		PageNum:    pageNum,
		TotalPages: totalPages,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
