package models

import "time"

type Session struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	TotalMarks  int       `db:"total_marks"`
	SessionDate time.Time `db:"session_date"`
	User        *User
}
