package models

import "database/sql"

type Quiz struct {
	ID           int            `db:"id"`
	SessionID    string         `db:"session_id"`
	QuestionID   string         `db:"question_id"`
	Answer       sql.NullString `db:"answer"`
	IsCorrect    sql.NullBool   `db:"isCorrect"`
	QuestionData *Question      // To hold the full question details for reports
}
