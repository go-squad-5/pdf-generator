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
	smtpPort = 1026
)

type EmailService struct {
	sessionRepo *repository.SessionRepository
	attemptRepo *repository.QuizAttemptRepository
	dialer      *gomail.Dialer
}

func NewEmailService(sessionRepo *repository.SessionRepository, attemptRepo *repository.QuizAttemptRepository) *EmailService {
	d := gomail.NewDialer(smtpHost, smtpPort, "", "")
	return &EmailService{
		sessionRepo: sessionRepo,
		attemptRepo: attemptRepo,
		dialer:      d,
	}
}

func (s *EmailService) SendQuizReportByEmail(sessionID int) error {
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
	wg.Wait()

	if sessionErr != nil {
		return fmt.Errorf("failed to fetch session data: %w", sessionErr)
	}
	if attemptsErr != nil {
		return fmt.Errorf("failed to fetch attempts data: %w", attemptsErr)
	}
	if session == nil {
		return fmt.Errorf("session with ID %d not found", sessionID)
	}
	if len(attempts) == 0 {
		return fmt.Errorf("no attempts found for session %d", sessionID)
	}

	const questionsPerEmail = 10
	var emailWg sync.WaitGroup

	for i := 0; i < len(attempts); i += questionsPerEmail {
		end := i + questionsPerEmail
		if end > len(attempts) {
			end = len(attempts)
		}
		paginatedAttempts := attempts[i:end]
		pageNumber := (i / questionsPerEmail) + 1
		totalPages := (len(attempts) + questionsPerEmail - 1) / questionsPerEmail

		emailWg.Add(1)
		go func(pAttempts []models.QuizAttempt, pNum, tPages int) {
			defer emailWg.Done()
			log.Printf("Goroutine started for email part %d/%d for session %d", pNum, tPages, sessionID)
			s.sendSingleEmailPart(session, pAttempts, pNum, tPages)
		}(paginatedAttempts, pageNumber, totalPages)
	}

	emailWg.Wait()
	log.Printf("All email parts for session %d have been processed.", sessionID)
	return nil
}

func (s *EmailService) sendSingleEmailPart(session *models.Session, attemptsChunk []models.QuizAttempt, pageNum, totalPages int) {
	body, err := s.parseEmailTemplate(session, attemptsChunk, pageNum, totalPages)
	if err != nil {
		log.Printf("ERROR: Could not parse email template for session %d part %d: %v", session.ID, pageNum, err)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "quiz-system@university.com")
	m.SetHeader("To", session.User.Email)
	m.SetHeader("Subject", fmt.Sprintf("Detailed Quiz Report for Session #%d (Part %d/%d)", session.ID, pageNum, totalPages))
	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("ERROR: Failed to send email for session %d part %d: %v", session.ID, pageNum, err)
	} else {
		log.Printf("SUCCESS: Sent email for session %d part %d to %s", session.ID, pageNum, session.User.Email)
	}
}

func (s *EmailService) parseEmailTemplate(session *models.Session, attempts []models.QuizAttempt, pageNum, totalPages int) (string, error) {
	const templateStr = `
    <!DOCTYPE html>
    <html>
    <head>
        <style>
            body { font-family: Arial, sans-serif; margin: 20px; }
            .container { border: 1px solid #ddd; padding: 20px; max-width: 600px; }
            h1, h2 { color: #333; }
            .question-block { margin-bottom: 15px; padding-bottom: 10px; border-bottom: 1px solid #eee; }
            .correct { color: green; }
            .incorrect { color: red; }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Detailed Quiz Report</h1>
            <h2>Session #{{.Session.ID}} - Part {{.PageNum}} of {{.TotalPages}}</h2>
            <p>Dear {{.Session.User.FirstName}}, here are the details for this part of your quiz:</p>

            {{range .Attempts}}
            <div class="question-block">
                <p><b>Question:</b> {{.Question.QuestionText}}</p>
                <p>Your Answer: <b class="{{if eq .ChosenOption .Question.CorrectOption}}correct{{else}}incorrect{{end}}">{{.ChosenOption}}</b></p>
                {{if ne .ChosenOption .Question.CorrectOption}}
                <p>Correct Answer: <b class="correct">{{.Question.CorrectOption}}</b></p>
                {{end}}
            </div>
            {{end}}
        </div>
    </body>
    </html>`

	t, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", err
	}

	data := struct {
		Session    *models.Session
		Attempts   []models.QuizAttempt
		PageNum    int
		TotalPages int
	}{
		Session:    session,
		Attempts:   attempts,
		PageNum:    pageNum,
		TotalPages: totalPages,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
