package repository

import (
  "database/sql"
  "regexp"
  "testing"

  "github.com/DATA-DOG/go-sqlmock"
  "github.com/stretchr/testify/assert"
)

func TestQuizzesRepository_GetQuizzesBySessionID(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
  }
  defer db.Close()

  repo := NewQuizzesRepository(db)

  t.Run("success", func(t *testing.T) {
    sessionID := "session-123"

    rows := sqlmock.NewRows([]string{"id", "session_id", "question_id", "answer", "is_correct", "question", "options", "correctAnswer", "topic"}).
      AddRow(1, sessionID, "q1", "A", true, "What is Go?", []byte(`["A", "B"]`), "A", "Programming").
      AddRow(2, sessionID, "q2", "C", false, "What is a slice?", []byte(`["C", "D"]`), "D", "Programming")

    query := `
		SELECT
			q.id, q.session_id, q.question_id, q.answer, q.is_correct,
			qs.question, qs.options, qs.answer as correctAnswer, qs.topic
		FROM quizzes q
		JOIN questions qs ON q.question_id = qs.id
		WHERE q.session_id = ?
		ORDER BY qs.id`

    mock.ExpectQuery(regexp.QuoteMeta(query)).
      WithArgs(sessionID).
      WillReturnRows(rows)

    quizzes, err := repo.GetQuizzesBySessionID(sessionID)

    assert.NoError(t, err)
    assert.NotNil(t, quizzes)
    assert.Len(t, quizzes, 2)
    assert.Equal(t, "What is Go?", quizzes[0].QuestionData.Question)
    assert.Equal(t, "What is a slice?", quizzes[1].QuestionData.Question)
  })

  t.Run("database error", func(t *testing.T) {
    sessionID := "session-error"

    query := `
		SELECT
			q.id, q.session_id, q.question_id, q.answer, q.is_correct,
			qs.question, qs.options, qs.answer as correctAnswer, qs.topic
		FROM quizzes q
		JOIN questions qs ON q.question_id = qs.id
		WHERE q.session_id = ?
		ORDER BY qs.id`

    mock.ExpectQuery(regexp.QuoteMeta(query)).
      WithArgs(sessionID).
      WillReturnError(sql.ErrConnDone)

    quizzes, err := repo.GetQuizzesBySessionID(sessionID)

    assert.Error(t, err)
    assert.Nil(t, quizzes)
  })

  t.Run("row-scan-error", func(t *testing.T) {
    sessionID := "session-123"

    rows := sqlmock.NewRows([]string{"id", "session_id", "question_id", "answer", "is_correct", "question", "options", "correctAnswer", "topic"}).
      AddRow(1, sessionID, "q1", "A", true, "What is Go?", []byte(`["A", "B"]`), "A", "Programming").
      AddRow(2, sessionID, "q2", "C", false, "What is a slice?", []byte(`["C", "D"]`), "D", "Programming").
      AddRow(2, sessionID, "q2", "C", false, "What is a slice?", []byte(`["C", "D"]`), "D", nil)

    query := `
		SELECT
			q.id, q.session_id, q.question_id, q.answer, q.is_correct,
			qs.question, qs.options, qs.answer as correctAnswer, qs.topic
		FROM quizzes q
		JOIN questions qs ON q.question_id = qs.id
		WHERE q.session_id = ?
		ORDER BY qs.id`

    mock.ExpectQuery(regexp.QuoteMeta(query)).
      WithArgs(sessionID).
      WillReturnRows(rows)

    quizzes, err := repo.GetQuizzesBySessionID(sessionID)

    assert.Error(t, err, "An error to be returned when a row scan fails")
    assert.Nil(t, quizzes, "The quizzes slice nil for rowscan error")
  })
}
