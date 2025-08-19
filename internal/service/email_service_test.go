package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// type mockDialer struct {
// 	sendCalled bool
// }

// func (m *mockDialer) DialAndSend(msgs ...*gomail.Message) error {
// 	m.sendCalled = true
// 	return nil
// }

func TestEmailService_SendQuizReportByEmail(t *testing.T) {
	t.Run("session not found", func(t *testing.T) {
		sessionRepo := &mockSessionRepo{session: nil}
		quizzesRepo := &mockQuizzesRepo{}

		service := NewEmailService(sessionRepo, quizzesRepo)
		err := service.SendQuizReportByEmail("sid-not-found")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "session with ID sid-not-found not found")
	})

	t.Run("repository error", func(t *testing.T) {
		sessionRepo := &mockSessionRepo{err: errors.New("db error")}
		quizzesRepo := &mockQuizzesRepo{}

		service := NewEmailService(sessionRepo, quizzesRepo)
		err := service.SendQuizReportByEmail("sid-error")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch data")
	})
}
