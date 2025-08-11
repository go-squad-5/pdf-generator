package repository

import (
	"database/sql"

	"github.com/go-squad-5/pdf-generator/internal/models"
)

type MarkRepository struct {
	DB *sql.DB
}

func NewMarkRepository(db *sql.DB) *MarkRepository {
	return &MarkRepository{DB: db}
}

func (r *MarkRepository) GetMarksByUserID(userID int) ([]models.Mark, error) {
	query := "SELECT subject, score FROM marks WHERE user_id = ?"
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var marks []models.Mark
	for rows.Next() {
		var mark models.Mark
		if err := rows.Scan(&mark.Subject, &mark.Score); err != nil {
			return nil, err
		}
		marks = append(marks, mark)
	}

	return marks, nil
}
