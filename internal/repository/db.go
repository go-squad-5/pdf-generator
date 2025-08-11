package repository

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA foreign_keys = ON;")

	createUsersTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"first_name" TEXT,
		"last_name" TEXT,
		"email" TEXT,
		"job_title" TEXT
	);`
	if _, err = db.Exec(createUsersTableSQL); err != nil {
		return nil, err
	}

	createMarksTableSQL := `CREATE TABLE IF NOT EXISTS marks (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"user_id" INTEGER,
		"subject" TEXT,
		"score" INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
	if _, err = db.Exec(createMarksTableSQL); err != nil {
		return nil, err
	}

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users")
	if err := row.Scan(&count); err != nil {
		return nil, err
	}

	if count == 0 {
		log.Println("Database is empty. Seeding with 500 new Indian sample data...")
		seedData(db)
	}

	return db, nil
}

func seedData(db *sql.DB) {
	rand.Seed(time.Now().UnixNano())

	firstNames := []string{"Aarav", "Vivaan", "Aditya", "Vihaan", "Arjun", "Sai", "Reyansh", "Ayaan", "Krishna", "Ishaan", "Priya", "Riya", "Saanvi", "Ananya", "Aadhya", "Ira", "Diya", "Avni", "Gauri", "Anika"}
	lastNames := []string{"Sharma", "Verma", "Gupta", "Singh", "Patel", "Kumar", "Das", "Mehta", "Jain", "Reddy", "Yadav", "Mishra", "Chauhan", "Malhotra", "Kapoor"}

	insertUserSQL := `INSERT INTO users(first_name, last_name, email, job_title) VALUES (?, ?, ?, ?)`
	userStmt, err := db.Prepare(insertUserSQL)
	if err != nil {
		log.Fatalf("Failed to prepare user seed statement: %v", err)
	}
	defer userStmt.Close()

	log.Println("Starting to seed 500 users. This may take a moment...")
	for i := 0; i < 500; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		email := fmt.Sprintf("%s.%s.%d@example.com", firstName, lastName, i)
		user := struct {
			FirstName string
			LastName  string
			Email     string
			JobTitle  string
		}{
			firstName,
			lastName,
			email,
			"Student",
		}

		res, err := userStmt.Exec(user.FirstName, user.LastName, user.Email, user.JobTitle)
		if err != nil {
			log.Printf("Failed to insert user seed data: %v", err)
		}

		userID, _ := res.LastInsertId()
		seedMarksForUser(db, int(userID), i)
	}

	log.Println("Database seeded successfully with 500 users.")
}

func seedMarksForUser(db *sql.DB, userID, userIndex int) {
	marksData := [][][]interface{}{
		{{"Physics", 88}, {"Chemistry", 92}, {"Biology", 85}, {"Mathematics", 79}, {"Computer Science", 81}, {"Hindi", 85}},
		{{"Physics", 72}, {"Chemistry", 68}, {"Biology", 75}, {"Mathematics", 81}, {"Computer Science", 77}, {"Hindi", 78}},
		{{"Physics", 94}, {"Chemistry", 90}, {"Biology", 88}, {"Mathematics", 85}, {"Computer Science", 91}, {"Hindi", 86}},
		{{"Physics", 65}, {"Chemistry", 71}, {"Biology", 62}, {"Mathematics", 70}, {"Computer Science", 69}, {"Hindi", 73}},
		{{"Physics", 80}, {"Chemistry", 85}, {"Biology", 78}, {"Mathematics", 88}, {"Computer Science", 79}, {"Hindi", 83}},
	}

	insertMarkSQL := `INSERT INTO marks(user_id, subject, score) VALUES (?, ?, ?)`
	markStmt, err := db.Prepare(insertMarkSQL)
	if err != nil {
		log.Fatalf("Failed to prepare mark seed statement: %v", err)
	}
	defer markStmt.Close()

	markSet := marksData[userIndex%len(marksData)]
	for _, mark := range markSet {
		_, err := markStmt.Exec(userID, mark[0], mark[1])
		if err != nil {
			log.Printf("Failed to insert mark seed data: %v", err)
		}
	}
}
