package repository

import (
	"database/sql"

	"github.com/go-squad-5/pdf-generator/internal/models"
)

type SessionRepository struct{ DB *sql.DB }

func NewSessionRepository(db *sql.DB) *SessionRepository { return &SessionRepository{DB: db} }

func (r *SessionRepository) GetSessionByID(sessionID int) (*models.Session, error) {
	query := `
		SELECT s.id, s.user_id, s.total_marks, s.session_date, u.first_name, u.last_name, u.email
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.id = ?`
	row := r.DB.QueryRow(query, sessionID)

	var session models.Session
	var user models.User
	err := row.Scan(&session.ID, &session.UserID, &session.TotalMarks, &session.SessionDate, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	session.User = &user
	return &session, nil
}
