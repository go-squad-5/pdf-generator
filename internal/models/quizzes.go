package models

type Quiz struct {
	ID           int       `db:"id"`
	SessionID    string    `db:"session_id"`
	QuestionID   string    `db:"question_id"`
	Answer       string    `db:"answer"`
	IsCorrect    bool      `db:"isCorrect"`
	QuestionData *Question // To hold the full question details for reports
}
