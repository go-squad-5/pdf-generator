package repository

import (
	"database/sql"

	"github.com/go-squad-5/pdf-generator/internal/models"
)

type QuizAttemptRepository struct{ DB *sql.DB }

func NewQuizAttemptRepository(db *sql.DB) *QuizAttemptRepository {
	return &QuizAttemptRepository{DB: db}
}

func (r *QuizAttemptRepository) GetAttemptsBySessionID(sessionID int) ([]models.QuizAttempt, error) {
	query := `
		SELECT
			qa.chosen_option,
			q.id, q.question_text, q.option_a, q.option_b, q.option_c, q.option_d, q.correct_option
		FROM quiz_attempts qa
		JOIN questions q ON qa.question_id = q.id
		WHERE qa.session_id = ?
		ORDER BY q.id ASC`
	rows, err := r.DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []models.QuizAttempt
	for rows.Next() {
		var attempt models.QuizAttempt
		var question models.Question
		err := rows.Scan(
			&attempt.ChosenOption,
			&question.ID, &question.QuestionText, &question.OptionA, &question.OptionB, &question.OptionC, &question.OptionD, &question.CorrectOption,
		)
		if err != nil {
			return nil, err
		}
		attempt.Question = &question
		attempts = append(attempts, attempt)
	}
	return attempts, nil
}
