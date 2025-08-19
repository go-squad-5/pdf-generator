package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-squad-5/pdf-generator/internal/models"
)

type MockQuizzesRepository struct{}

func makeQuiz(id int, ssid, qid string, ans string, isCorrect bool) models.Quiz {
	return models.Quiz{
		ID:         id,
		SessionID:  ssid,
		QuestionID: qid,
		Answer:     sql.NullString{String: ans, Valid: true},
		IsCorrect:  sql.NullBool{Bool: isCorrect, Valid: true},
		QuestionData: &models.Question{
			Question: fmt.Sprintf("Question %d?", id),
			Options: models.OptionsMap{
				1: "Option A",
				2: "Option B",
				3: "Option C",
				4: "Option D",
			},
			Answer: "a",
			Topic:  "Sample",
		},
	}
}

func (m *MockQuizzesRepository) GetQuizzesBySessionID(sessionID string) ([]models.Quiz, error) {
	var quizzes []models.Quiz
	for i := 0; i < 10; i++ {
		quizzes = append(quizzes, makeQuiz(i, fmt.Sprintf("session-%d", i), fmt.Sprintf("question-%d", i), fmt.Sprintf("option-%d", i%4+1), true))
	}
	log.Print(quizzes)
	return quizzes, nil
}
