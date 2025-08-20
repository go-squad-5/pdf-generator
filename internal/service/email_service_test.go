package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/go-squad-5/pdf-generator/internal/models"
	"github.com/stretchr/testify/assert"
)

//type mockDialer struct {
//	sendCalled bool
//}
//
//func (m *mockDialer) DialAndSend(msgs ...*gomail.Message) error {
//	m.sendCalled = true
//	return nil
//}

func TestEmailService_SendQuizReportByEmail(t *testing.T) {
	t.Run("session not found", func(t *testing.T) {
		sessionRepo := &mockSessionRepo{session: nil}
		quizzesRepo := &mockQuizzesRepo{}

		service := NewEmailService(sessionRepo, quizzesRepo)
		err := service.SendQuizReportByEmail("ssid-not-found")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "session with ID ssid-not-found not found")
	})

	t.Run("repository error", func(t *testing.T) {
		sessionRepo := &mockSessionRepo{err: errors.New("db error")}
		quizzesRepo := &mockQuizzesRepo{}

		service := NewEmailService(sessionRepo, quizzesRepo)
		err := service.SendQuizReportByEmail("sid-error")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch data")
	})

	t.Run("No-error", func(t *testing.T) {
		sessionRepo := &mockSessionRepo{
			session: &models.Session{SessionID: "valid-ssid", Email: "test@example.com"},
			err:     nil,
		}
		quizzesRepo := &mockQuizzesRepo{
			quizzes: []models.Quiz{
				{
					ID: 1,
					QuestionData: &models.Question{
						Question: "Is this test working now?",
						Answer:   "Yes",
					},
					Answer:    sql.NullString{String: "Yes", Valid: true},
					IsCorrect: sql.NullBool{Bool: true, Valid: true},
				},
			},
		}

		service := NewEmailService(sessionRepo, quizzesRepo)
		err := service.SendQuizReportByEmail("valid-ssid")

		assert.Nil(t, err)
	})

}
