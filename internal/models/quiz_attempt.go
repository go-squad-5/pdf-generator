package models

type QuizAttempt struct {
	ChosenOption string `db:"chosen_option"`
	Question     *Question
}
