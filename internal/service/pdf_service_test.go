package service

import (
	"errors"
	"testing"

	"github.com/go-squad-5/pdf-generator/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockSessionRepo struct {
	session *models.Session
	err     error
}

func (m *mockSessionRepo) GetSessionByID(sessionID string) (*models.Session, error) {
	return m.session, m.err
}

type mockQuizzesRepo struct {
	quizzes []models.Quiz
	err     error
}

func (m *mockQuizzesRepo) GetQuizzesBySessionID(sessionID string) ([]models.Quiz, error) {
	return m.quizzes, m.err
}

func TestPDFService_GenerateQuizReport(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sessionRepo := &mockSessionRepo{
			session: &models.Session{SessionID: "sid-1", Email: "test@test.com", Score: 1},
		}
		quizzesRepo := &mockQuizzesRepo{
			quizzes: []models.Quiz{{ID: 1, QuestionData: &models.Question{Question: "Q1"}}},
		}

		service := NewPDFService(sessionRepo, quizzesRepo)
		pdfBytes, err := service.GenerateQuizReport("sid-1")

		assert.NoError(t, err)
		assert.NotEmpty(t, pdfBytes)
	})

	t.Run("session not found", func(t *testing.T) {
		sessionRepo := &mockSessionRepo{
			session: nil,
			err:     nil,
		}
		quizzesRepo := &mockQuizzesRepo{}

		service := NewPDFService(sessionRepo, quizzesRepo)
		pdfBytes, err := service.GenerateQuizReport("sid-not-found")

		assert.Error(t, err)
		assert.Nil(t, pdfBytes)
		assert.Contains(t, err.Error(), "session with ID sid-not-found not found")
	})

	t.Run("repository error", func(t *testing.T) {
		sessionRepo := &mockSessionRepo{
			err: errors.New("db error"),
		}
		quizzesRepo := &mockQuizzesRepo{}

		service := NewPDFService(sessionRepo, quizzesRepo)
		pdfBytes, err := service.GenerateQuizReport("sid-error")

		assert.Error(t, err)
		assert.Nil(t, pdfBytes)
		assert.Contains(t, err.Error(), "failed to fetch data")
	})
}
