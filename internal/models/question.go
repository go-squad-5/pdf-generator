package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// OptionsMap is a custom type to handle the JSON 'options' field.
type OptionsMap map[string]string

// Value converts the OptionsMap to a JSON string for database storage.
func (o OptionsMap) Value() (driver.Value, error) {
	return json.Marshal(o)
}

// Scan converts the JSON string from the database into an OptionsMap.
func (o *OptionsMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &o)
}

// Question matches your new schema.
type Question struct {
	ID       string     `db:"id"`
	Question string     `db:"question"`
	Options  OptionsMap `db:"options"` // Using our custom JSON type
	Answer   string     `db:"answer"`
	Topic    string     `db:"topic"`
}
