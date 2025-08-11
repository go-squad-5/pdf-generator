package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-squad-5/pdf-generator/internal/models"
	"github.com/google/uuid"
)

func InitDB() (*sql.DB, error) {
	dsn := "root:root@tcp(127.0.0.1:3333)/quizdb?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not connect to MySQL: %w", err)
	}

	return db, nil
}

func SeedData(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Session").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("Database already contains data. Seeding skipped.")
		return nil
	}

	rand.Seed(time.Now().UnixNano())

	log.Println("Seeding Questions...")
	questionIDs := seedQuestions(db)

	log.Println("Seeding Sessions and Quizzes...")
	seedSessionsAndQuizzes(db, questionIDs)

	return nil
}

func seedQuestions(db *sql.DB) []string {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO Questions(id, question, options, answer, topic) VALUES (?, ?, ?, ?, ?)`)
	defer stmt.Close()

	var questionIDs []string
	questions := getFullQuestionList()

	for _, q := range questions {
		questionID := uuid.New().String()
		optionsJSON, _ := json.Marshal(q.Options)
		_, err := stmt.Exec(questionID, q.Text, optionsJSON, q.Correct, q.Topic)
		if err != nil {
			log.Printf("Failed to insert question: %v", err)
		}
		questionIDs = append(questionIDs, questionID)
	}
	tx.Commit()
	return questionIDs
}

func seedSessionsAndQuizzes(db *sql.DB, questionIDs []string) {
	type questionInfo struct {
		correctAnswer string
		optionsJSON   string
	}
	questionsMap := make(map[string]questionInfo)
	rows, err := db.Query("SELECT id, answer, options FROM Questions")
	if err != nil {
		log.Fatalf("Failed to pre-fetch questions: %v", err)
	}
	for rows.Next() {
		var id string
		var info questionInfo
		rows.Scan(&id, &info.correctAnswer, &info.optionsJSON)
		questionsMap[id] = info
	}
	rows.Close()

	sessionTx, _ := db.Begin()
	sessionStmt, _ := sessionTx.Prepare(`INSERT INTO Session(session_id, email, topic, score) VALUES (?, ?, ?, ?)`)
	defer sessionStmt.Close()

	quizTx, _ := db.Begin()
	quizStmt, _ := quizTx.Prepare(`INSERT INTO Quizzes(session_id, question_id, answer, isCorrect) VALUES (?, ?, ?, ?)`)
	defer quizStmt.Close()

	users := []string{"priya.sharma@example.com", "rohan.verma@example.com", "anjali.singh@example.com"}

	for _, email := range users {
		for i := 0; i < 3; i++ {
			sessionID := uuid.New().String()
			topic := "General Knowledge"
			score := 0

			rand.Shuffle(len(questionIDs), func(i, j int) {
				questionIDs[i], questionIDs[j] = questionIDs[j], questionIDs[i]
			})
			quizQuestionIDs := questionIDs[:25]

			for _, qID := range quizQuestionIDs {
				qInfo := questionsMap[qID]

				var options map[string]string
				json.Unmarshal([]byte(qInfo.optionsJSON), &options)

				possibleAnswers := make([]string, 0, len(options))
				for k := range options {
					possibleAnswers = append(possibleAnswers, k)
				}

				chosenAnswer := possibleAnswers[rand.Intn(len(possibleAnswers))]
				isCorrect := chosenAnswer == qInfo.correctAnswer

				if isCorrect {
					score++
				}
				quizStmt.Exec(sessionID, qID, chosenAnswer, isCorrect)
			}
			sessionStmt.Exec(sessionID, email, topic, score)
		}
	}
	sessionTx.Commit()
	quizTx.Commit()
}

func getFullQuestionList() []struct {
	Text, Topic, Correct string
	Options              models.OptionsMap
} {
	return []struct {
		Text, Topic, Correct string
		Options              models.OptionsMap
	}{
		{"What is the capital of France?", "Geography", "Paris", models.OptionsMap{"Berlin": "Berlin", "Madrid": "Madrid", "Paris": "Paris", "Rome": "Rome"}},
		{"Which planet is known as the Red Planet?", "Science", "Mars", models.OptionsMap{"Earth": "Earth", "Mars": "Mars", "Jupiter": "Jupiter", "Venus": "Venus"}},
		{"What is the largest ocean on Earth?", "Geography", "Pacific", models.OptionsMap{"Atlantic": "Atlantic", "Indian": "Indian", "Arctic": "Arctic", "Pacific": "Pacific"}},
		{"Who wrote 'To Kill a Mockingbird'?", "Literature", "Harper Lee", models.OptionsMap{"Harper Lee": "Harper Lee", "Mark Twain": "Mark Twain", "J.K. Rowling": "J.K. Rowling", "F. Scott Fitzgerald": "F. Scott Fitzgerald"}},
		{"What is the chemical symbol for water?", "Science", "H2O", models.OptionsMap{"O2": "O2", "H2O": "H2O", "CO2": "CO2", "NaCl": "NaCl"}},
		{"In which year did the Titanic sink?", "History", "1912", models.OptionsMap{"1905": "1905", "1912": "1912", "1918": "1918", "1923": "1923"}},
		{"What is the currency of Japan?", "World", "Yen", models.OptionsMap{"Won": "Won", "Yuan": "Yuan", "Yen": "Yen", "Dollar": "Dollar"}},
		{"Who painted the Mona Lisa?", "Art", "Leonardo da Vinci", models.OptionsMap{"Vincent van Gogh": "Vincent van Gogh", "Pablo Picasso": "Pablo Picasso", "Leonardo da Vinci": "Leonardo da Vinci", "Claude Monet": "Claude Monet"}},
		{"What is the hardest natural substance on Earth?", "Science", "Diamond", models.OptionsMap{"Gold": "Gold", "Iron": "Iron", "Diamond": "Diamond", "Platinum": "Platinum"}},
		{"Which element has the atomic number 1?", "Science", "Hydrogen", models.OptionsMap{"Helium": "Helium", "Oxygen": "Oxygen", "Hydrogen": "Hydrogen", "Carbon": "Carbon"}},
		{"What is the capital of Australia?", "Geography", "Canberra", models.OptionsMap{"Sydney": "Sydney", "Melbourne": "Melbourne", "Canberra": "Canberra", "Perth": "Perth"}},
		{"Who discovered penicillin?", "History", "Alexander Fleming", models.OptionsMap{"Marie Curie": "Marie Curie", "Albert Einstein": "Albert Einstein", "Isaac Newton": "Isaac Newton", "Alexander Fleming": "Alexander Fleming"}},
		{"What is the tallest mountain in the world?", "Geography", "Mount Everest", models.OptionsMap{"K2": "K2", "Kangchenjunga": "Kangchenjunga", "Mount Everest": "Mount Everest", "Lhotse": "Lhotse"}},
		{"Which country is known as the Land of the Rising Sun?", "World", "Japan", models.OptionsMap{"China": "China", "South Korea": "South Korea", "Japan": "Japan", "Thailand": "Thailand"}},
		{"What is the main ingredient in guacamole?", "Food", "Avocado", models.OptionsMap{"Tomato": "Tomato", "Avocado": "Avocado", "Onion": "Onion", "Lime": "Lime"}},
		{"How many continents are there?", "Geography", "7", models.OptionsMap{"5": "5", "6": "6", "7": "7", "8": "8"}},
		{"Who was the first person to walk on the moon?", "History", "Neil Armstrong", models.OptionsMap{"Buzz Aldrin": "Buzz Aldrin", "Yuri Gagarin": "Yuri Gagarin", "Michael Collins": "Michael Collins", "Neil Armstrong": "Neil Armstrong"}},
		{"What is the largest desert in the world?", "Geography", "Antarctic Polar Desert", models.OptionsMap{"Sahara Desert": "Sahara Desert", "Arabian Desert": "Arabian Desert", "Gobi Desert": "Gobi Desert", "Antarctic Polar Desert": "Antarctic Polar Desert"}},
		{"Which is the longest river in the world?", "Geography", "Nile River", models.OptionsMap{"Amazon River": "Amazon River", "Nile River": "Nile River", "Yangtze River": "Yangtze River", "Mississippi River": "Mississippi River"}},
		{"What does 'CPU' stand for?", "Technology", "Central Processing Unit", models.OptionsMap{"Central Process Unit": "Central Process Unit", "Computer Personal Unit": "Computer Personal Unit", "Central Processing Unit": "Central Processing Unit", "Computer Primary Unit": "Computer Primary Unit"}},
		{"Who wrote the play 'Romeo and Juliet'?", "Literature", "William Shakespeare", models.OptionsMap{"Charles Dickens": "Charles Dickens", "William Shakespeare": "William Shakespeare", "George Orwell": "George Orwell", "Jane Austen": "Jane Austen"}},
		{"What is the chemical symbol for gold?", "Science", "Au", models.OptionsMap{"Ag": "Ag", "Au": "Au", "Pb": "Pb", "Fe": "Fe"}},
		{"Which is the smallest planet in our solar system?", "Science", "Mercury", models.OptionsMap{"Venus": "Venus", "Mars": "Mars", "Mercury": "Mercury", "Uranus": "Uranus"}},
		{"What is the capital of Canada?", "Geography", "Ottawa", models.OptionsMap{"Toronto": "Toronto", "Vancouver": "Vancouver", "Montreal": "Montreal", "Ottawa": "Ottawa"}},
		{"How many bones are in the adult human body?", "Science", "206", models.OptionsMap{"206": "206", "208": "208", "210": "210", "212": "212"}},
		{"Which artist is known for the 'Starry Night' painting?", "Art", "Vincent van Gogh", models.OptionsMap{"Pablo Picasso": "Pablo Picasso", "Claude Monet": "Claude Monet", "Salvador Dalí": "Salvador Dalí", "Vincent van Gogh": "Vincent van Gogh"}},
		{"What is the primary language spoken in Brazil?", "World", "Portuguese", models.OptionsMap{"Spanish": "Spanish", "Portuguese": "Portuguese", "Brazilian": "Brazilian", "English": "English"}},
		{"Who invented the telephone?", "History", "Alexander Graham Bell", models.OptionsMap{"Thomas Edison": "Thomas Edison", "Nikola Tesla": "Nikola Tesla", "Alexander Graham Bell": "Alexander Graham Bell", "Guglielmo Marconi": "Guglielmo Marconi"}},
		{"What is the capital of Egypt?", "Geography", "Cairo", models.OptionsMap{"Alexandria": "Alexandria", "Giza": "Giza", "Cairo": "Cairo", "Luxor": "Luxor"}},
		{"Which gas do plants absorb from the atmosphere?", "Science", "Carbon Dioxide", models.OptionsMap{"Oxygen": "Oxygen", "Nitrogen": "Nitrogen", "Carbon Dioxide": "Carbon Dioxide", "Hydrogen": "Hydrogen"}},
		{"What is the main component of the sun?", "Science", "Hydrogen and Helium", models.OptionsMap{"Liquid lava": "Liquid lava", "Rock": "Rock", "Hydrogen and Helium": "Hydrogen and Helium", "Oxygen": "Oxygen"}},
		{"Who was the first female Prime Minister of the United Kingdom?", "History", "Margaret Thatcher", models.OptionsMap{"Theresa May": "Theresa May", "Margaret Thatcher": "Margaret Thatcher", "Angela Merkel": "Angela Merkel", "Indira Gandhi": "Indira Gandhi"}},
		{"What is the largest country by land area?", "Geography", "Russia", models.OptionsMap{"Canada": "Canada", "China": "China", "USA": "USA", "Russia": "Russia"}},
		{"Which of these is a primary color?", "Art", "Blue", models.OptionsMap{"Green": "Green", "Orange": "Orange", "Blue": "Blue", "Purple": "Purple"}},
		{"What is the boiling point of water at sea level?", "Science", "100°C", models.OptionsMap{"90°C": "90°C", "100°C": "100°C", "110°C": "110°C", "120°C": "120°C"}},
		{"Who is the author of the Harry Potter series?", "Literature", "J.K. Rowling", models.OptionsMap{"J.R.R. Tolkien": "J.R.R. Tolkien", "George R.R. Martin": "George R.R. Martin", "Suzanne Collins": "Suzanne Collins", "J.K. Rowling": "J.K. Rowling"}},
		{"What is the capital of Italy?", "Geography", "Rome", models.OptionsMap{"Milan": "Milan", "Naples": "Naples", "Rome": "Rome", "Venice": "Venice"}},
		{"Which ocean is the Bermuda Triangle located in?", "Geography", "Atlantic", models.OptionsMap{"Atlantic": "Atlantic", "Pacific": "Pacific", "Indian": "Indian", "Arctic": "Arctic"}},
		{"What is the square root of 64?", "Math", "8", models.OptionsMap{"6": "6", "7": "7", "8": "8", "9": "9"}},
		{"What is the national sport of Canada?", "Sports", "Hockey and Lacrosse", models.OptionsMap{"Hockey and Lacrosse": "Hockey and Lacrosse", "Baseball": "Baseball", "Basketball": "Basketball", "Soccer": "Soccer"}},
		{"Which composer was deaf for the last part of his life?", "Music", "Beethoven", models.OptionsMap{"Mozart": "Mozart", "Bach": "Bach", "Beethoven": "Beethoven", "Chopin": "Chopin"}},
		{"What is the capital of Spain?", "Geography", "Madrid", models.OptionsMap{"Barcelona": "Barcelona", "Seville": "Seville", "Lisbon": "Lisbon", "Madrid": "Madrid"}},
		{"What is the largest mammal in the world?", "Science", "Blue Whale", models.OptionsMap{"Elephant": "Elephant", "Blue Whale": "Blue Whale", "Giraffe": "Giraffe", "Great White Shark": "Great White Shark"}},
		{"Which country gifted the Statue of Liberty to the USA?", "History", "France", models.OptionsMap{"Germany": "Germany", "United Kingdom": "United Kingdom", "France": "France", "Italy": "Italy"}},
		{"What is the chemical formula for table salt?", "Science", "NaCl", models.OptionsMap{"H2O": "H2O", "C6H12O6": "C6H12O6", "NaCl": "NaCl", "CO2": "CO2"}},
		{"Who discovered gravity?", "History", "Isaac Newton", models.OptionsMap{"Albert Einstein": "Albert Einstein", "Galileo Galilei": "Galileo Galilei", "Isaac Newton": "Isaac Newton", "Nikola Tesla": "Nikola Tesla"}},
		{"What is the capital of Russia?", "Geography", "Moscow", models.OptionsMap{"Saint Petersburg": "Saint Petersburg", "Kazan": "Kazan", "Novosibirsk": "Novosibirsk", "Moscow": "Moscow"}},
		{"In what year did World War II end?", "History", "1945", models.OptionsMap{"1943": "1943", "1944": "1944", "1945": "1945", "1946": "1946"}},
		{"What is the most spoken language in the world?", "World", "Mandarin Chinese", models.OptionsMap{"English": "English", "Spanish": "Spanish", "Mandarin Chinese": "Mandarin Chinese", "Hindi": "Hindi"}},
		{"Who is known as the 'Father of Computers'?", "Technology", "Charles Babbage", models.OptionsMap{"Alan Turing": "Alan Turing", "Charles Babbage": "Charles Babbage", "Tim Berners-Lee": "Tim Berners-Lee", "Steve Jobs": "Steve Jobs"}},
	}
}
