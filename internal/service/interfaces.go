package service

import (
	"github.com/go-squad-5/pdf-generator/internal/models"
	"gopkg.in/gomail.v2"
)

type SessionRepo interface {
	GetSessionByID(sessionID string) (*models.Session, error)
}

type QuizzesRepo interface {
	GetQuizzesBySessionID(sessionID string) ([]models.Quiz, error)
}

type MailDialer interface {
	DialAndSend(m ...*gomail.Message) error
}
