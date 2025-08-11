package models

type Question struct {
	ID            int    `db:"id"`
	QuestionText  string `db:"question_text"`
	OptionA       string `db:"option_a"`
	OptionB       string `db:"option_b"`
	OptionC       string `db:"option_c"`
	OptionD       string `db:"option_d"`
	CorrectOption string `db:"correct_option"`
}
