package repository

// import (
// 	"database/sql"
// 	"regexp"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/go-squad-5/pdf-generator/internal/models"
// 	"github.com/stretchr/testify/assert"
// )

// func TestQuizzesRepository_GetQuizzesBySessionID(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	repo := NewQuizzesRepository(db)

// 	t.Run("success", func(t *testing.T) {
// 		sessionID := "session-123"
// 		expectedQuizzes := []models.Quiz{
// 			{
// 				ID:         1,
// 				SessionID:  sessionID,
// 				QuestionID: "q1",
// 				Answer:     sql.NullString{String: "A", Valid: true},
// 				IsCorrect:  sql.NullBool{Bool: true, Valid: true},
// 				QuestionData: &models.Question{
// 					Question: "What is Go?",
// 					Options:  []string{"A", "B"},
// 					Answer:   "A",
// 					Topic:    "Programming",
// 				},
// 			},
// 		}

// 		rows := sqlmock.NewRows([]string{"id", "session_id", "question_id", "answer", "is_correct", "question", "options", "correctAnswer", "topic"}).
// 			AddRow(1, sessionID, "q1", "A", true, "What is Go?", `["A", "B"]`, "A", "Programming")

// 		query := `
//         SELECT
//             q.id, q.session_id, q.question_id, q.answer, q.is_correct,
//             qs.question, qs.options, qs.answer as correctAnswer, qs.topic
//         FROM quizzes q
//         JOIN questions qs ON q.question_id = qs.id
//         WHERE q.session_id = ?
//         ORDER BY qs.id`
// 		mock.ExpectQuery(regexp.QuoteMeta(query)).
// 			WithArgs(sessionID).
// 			WillReturnRows(rows)

// 		quizzes, err := repo.GetQuizzesBySessionID(sessionID)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, quizzes)
// 		assert.Len(t, quizzes, 1)
// 		assert.Equal(t, expectedQuizzes[0].ID, quizzes[0].ID)
// 		assert.Equal(t, expectedQuizzes[0].QuestionData.Question, quizzes[0].QuestionData.Question)
// 	})

// 	t.Run("database error", func(t *testing.T) {
// 		sessionID := "session-error"
// 		query := `
//         SELECT
//             q.id, q.session_id, q.question_id, q.answer, q.is_correct,
//             qs.question, qs.options, qs.answer as correctAnswer, qs.topic
//         FROM quizzes q
//         JOIN questions qs ON q.question_id = qs.id
//         WHERE q.session_id = ?
//         ORDER BY qs.id`
// 		mock.ExpectQuery(regexp.QuoteMeta(query)).
// 			WithArgs(sessionID).
// 			WillReturnError(sql.ErrConnDone)

// 		quizzes, err := repo.GetQuizzesBySessionID(sessionID)

// 		assert.Error(t, err)
// 		assert.Nil(t, quizzes)
// 	})
// }
