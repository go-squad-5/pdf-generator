package models

type Session struct {
	SessionID string `db:"session_id"`
	Email     string `db:"email"`
	Topic     string `db:"topic"`
	Score     int    `db:"score"`
}
