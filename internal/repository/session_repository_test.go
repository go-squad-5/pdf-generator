package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-squad-5/pdf-generator/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSessionRepository_GetSessionByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewSessionRepository(db)

	t.Run("success", func(t *testing.T) {
		expectedSession := &models.Session{
			SessionID: "test-uuid",
			Email:     "test@example.com",
			Topic:     "Go",
			Score:     85,
		}

		rows := sqlmock.NewRows([]string{"session_id", "email", "topic", "score"}).
			AddRow(expectedSession.SessionID, expectedSession.Email, expectedSession.Topic, expectedSession.Score)

		query := "SELECT session_id, email, topic, score FROM SESSION_TABLE WHERE session_id = ?"
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("test-uuid").
			WillReturnRows(rows)

		session, err := repo.GetSessionByID("test-uuid")

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, expectedSession, session)
	})

	t.Run("not found", func(t *testing.T) {
		query := "SELECT session_id, email, topic, score FROM SESSION_TABLE WHERE session_id = ?"
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("not-found-uuid").
			WillReturnError(sql.ErrNoRows)

		session, err := repo.GetSessionByID("not-found-uuid")

		assert.NoError(t, err)
		assert.Nil(t, session)
	})

	t.Run("database error", func(t *testing.T) {
		query := "SELECT session_id, email, topic, score FROM SESSION_TABLE WHERE session_id = ?"
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("error-uuid").
			WillReturnError(sql.ErrConnDone)

		session, err := repo.GetSessionByID("error-uuid")

		assert.Error(t, err)
		assert.Nil(t, session)
	})
}
