package repository

import (
	"database/sql"

	"github.com/go-squad-5/pdf-generator/internal/models"
)

type QuizzesRepository struct{ DB *sql.DB }

func NewQuizzesRepository(db *sql.DB) *QuizzesRepository { return &QuizzesRepository{DB: db} }

func (r *QuizzesRepository) GetQuizzesBySessionID(sessionID string) ([]models.Quiz, error) {
	query := `
		SELECT
			q.id, q.session_id, q.question_id, q.answer, q.is_correct,
			qs.question, qs.options, qs.answer as correctAnswer, qs.topic
		FROM quizzes q
		JOIN questions qs ON q.question_id = qs.id
		WHERE q.session_id = ?
		ORDER BY qs.id`
	rows, err := r.DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []models.Quiz
	for rows.Next() {
		var quiz models.Quiz
		var question models.Question
		var correctAnswer string

		err := rows.Scan(
			&quiz.ID, &quiz.SessionID, &quiz.QuestionID, &quiz.Answer, &quiz.IsCorrect,
			&question.Question, &question.Options, &correctAnswer, &question.Topic,
		)
		if err != nil {
			return nil, err
		}
		question.Answer = correctAnswer
		quiz.QuestionData = &question
		quizzes = append(quizzes, quiz)
	}
	return quizzes, nil
}
