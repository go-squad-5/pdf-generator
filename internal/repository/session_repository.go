package repository

import (
	"database/sql"

	"github.com/go-squad-5/pdf-generator/internal/models"
)

type SessionRepository struct{ DB *sql.DB }

func NewSessionRepository(db *sql.DB) *SessionRepository { return &SessionRepository{DB: db} }

func (r *SessionRepository) GetSessionByID(sessionID string) (*models.Session, error) {
	query := "SELECT session_id, email, topic, score FROM Session WHERE session_id = ?"
	row := r.DB.QueryRow(query, sessionID)

	var session models.Session
	err := row.Scan(&session.SessionID, &session.Email, &session.Topic, &session.Score)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}
