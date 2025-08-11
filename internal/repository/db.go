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

	createUsersTableSQL := `CREATE TABLE IF NOT EXISTS users ( "id" INTEGER NOT NULL PRIMARY KEY, "first_name" TEXT, "last_name" TEXT, "email" TEXT, "job_title" TEXT );`
	createQuestionsTableSQL := `CREATE TABLE IF NOT EXISTS questions ( "id" INTEGER NOT NULL PRIMARY KEY, "question_text" TEXT, "option_a" TEXT, "option_b" TEXT, "option_c" TEXT, "option_d" TEXT, "correct_option" TEXT );`
	createSessionsTableSQL := `CREATE TABLE IF NOT EXISTS sessions ( "id" INTEGER NOT NULL PRIMARY KEY, "user_id" INTEGER, "total_marks" INTEGER, "session_date" DATETIME, FOREIGN KEY(user_id) REFERENCES users(id) );`
	createAttemptsTableSQL := `CREATE TABLE IF NOT EXISTS quiz_attempts ( "id" INTEGER NOT NULL PRIMARY KEY, "session_id" INTEGER, "question_id" INTEGER, "chosen_option" TEXT, FOREIGN KEY(session_id) REFERENCES sessions(id), FOREIGN KEY(question_id) REFERENCES questions(id) );`

	for _, query := range []string{createUsersTableSQL, createQuestionsTableSQL, createSessionsTableSQL, createAttemptsTableSQL} {
		if _, err = db.Exec(query); err != nil {
			return nil, err
		}
	}

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users")
	if err := row.Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		log.Println("Database is empty. Seeding with multiple users and sessions...")
		seedData(db)
	}

	return db, nil
}

func seedData(db *sql.DB) {
	rand.Seed(time.Now().UnixNano())
	log.Println("Seeding 100 questions...")
	questionIDs := seedQuestions(db)

	log.Println("Seeding 10 users...")
	userIDs := seedUsers(db)

	log.Println("Seeding multiple quiz sessions for each user...")
	for _, userID := range userIDs {
		for i := 0; i < 3; i++ {
			seedQuizForUser(db, userID, questionIDs)
		}
	}

	log.Println("Database seeding complete.")
}

func seedQuestions(db *sql.DB) []int64 {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO questions(question_text, option_a, option_b, option_c, option_d, correct_option) VALUES (?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	var questionIDs []int64
	questions := getFullQuestionList()
	for _, q := range questions {
		res, _ := stmt.Exec(q.Text, q.A, q.B, q.C, q.D, q.Correct)
		id, _ := res.LastInsertId()
		questionIDs = append(questionIDs, id)
	}
	tx.Commit()
	return questionIDs
}

func seedUsers(db *sql.DB) []int64 {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO users(first_name, last_name, email, job_title) VALUES (?, ?, ?, ?)`)
	defer stmt.Close()

	var userIDs []int64
	users := []struct{ F, L string }{
		{"Aarav", "Mehta"}, {"Isha", "Patel"}, {"Rohan", "Kumar"}, {"Priya", "Singh"},
		{"Vikram", "Jain"}, {"Anika", "Sharma"}, {"Arjun", "Verma"}, {"Diya", "Gupta"},
		{"Kabir", "Das"}, {"Saanvi", "Reddy"},
	}
	for _, u := range users {
		res, _ := stmt.Exec(u.F, u.L, fmt.Sprintf("%s.%s@example.com", u.F, u.L), "Student")
		id, _ := res.LastInsertId()
		userIDs = append(userIDs, id)
	}
	tx.Commit()
	return userIDs
}

func seedQuizForUser(db *sql.DB, userID int64, allQuestionIDs []int64) {
	sessionTx, _ := db.Begin()
	sessionStmt, _ := sessionTx.Prepare(`INSERT INTO sessions(user_id, total_marks, session_date) VALUES (?, ?, ?)`)
	defer sessionStmt.Close()

	rand.Shuffle(len(allQuestionIDs), func(i, j int) {
		allQuestionIDs[i], allQuestionIDs[j] = allQuestionIDs[j], allQuestionIDs[i]
	})
	quizQuestionIDs := allQuestionIDs[:50]

	totalMarks := 0
	chosenOptions := make(map[int64]string)

	options := []string{"a", "b", "c", "d"}
	for _, qID := range quizQuestionIDs {
		var correctOpt string
		db.QueryRow("SELECT correct_option FROM questions WHERE id = ?", qID).Scan(&correctOpt)

		var chosenOpt string
		if rand.Intn(10) > 2 {
			chosenOpt = correctOpt
			totalMarks++
		} else {
			for {
				chosenOpt = options[rand.Intn(4)]
				if chosenOpt != correctOpt {
					break
				}
			}
		}
		chosenOptions[qID] = chosenOpt
	}

	res, _ := sessionStmt.Exec(userID, totalMarks, time.Now().Add(-time.Hour*time.Duration(rand.Intn(72))))
	sessionID, _ := res.LastInsertId()
	sessionTx.Commit()

	attemptTx, _ := db.Begin()
	attemptStmt, _ := attemptTx.Prepare(`INSERT INTO quiz_attempts(session_id, question_id, chosen_option) VALUES (?, ?, ?)`)
	defer attemptStmt.Close()
	for qID, chosenOpt := range chosenOptions {
		attemptStmt.Exec(sessionID, qID, chosenOpt)
	}
	attemptTx.Commit()
}
func getFullQuestionList() []struct{ Text, A, B, C, D, Correct string } {
	return []struct{ Text, A, B, C, D, Correct string }{
		{"What is the capital of France?", "Berlin", "Madrid", "Paris", "Rome", "c"},
		{"Which planet is known as the Red Planet?", "Earth", "Mars", "Jupiter", "Venus", "b"},
		{"What is the largest ocean on Earth?", "Atlantic", "Indian", "Arctic", "Pacific", "d"},
		{"Who wrote 'To Kill a Mockingbird'?", "Harper Lee", "Mark Twain", "J.K. Rowling", "F. Scott Fitzgerald", "a"},
		{"What is the chemical symbol for water?", "O2", "H2O", "CO2", "NaCl", "b"},
		{"In which year did the Titanic sink?", "1905", "1912", "1918", "1923", "b"},
		{"What is the currency of Japan?", "Won", "Yuan", "Yen", "Dollar", "c"},
		{"Who painted the Mona Lisa?", "Vincent van Gogh", "Pablo Picasso", "Leonardo da Vinci", "Claude Monet", "c"},
		{"What is the hardest natural substance on Earth?", "Gold", "Iron", "Diamond", "Platinum", "c"},
		{"Which element has the atomic number 1?", "Helium", "Oxygen", "Hydrogen", "Carbon", "c"},
		{"What is the capital of Australia?", "Sydney", "Melbourne", "Canberra", "Perth", "c"},
		{"Who discovered penicillin?", "Marie Curie", "Albert Einstein", "Isaac Newton", "Alexander Fleming", "d"},
		{"What is the tallest mountain in the world?", "K2", "Kangchenjunga", "Mount Everest", "Lhotse", "c"},
		{"Which country is known as the Land of the Rising Sun?", "China", "South Korea", "Japan", "Thailand", "c"},
		{"What is the main ingredient in guacamole?", "Tomato", "Avocado", "Onion", "Lime", "b"},
		{"How many continents are there?", "5", "6", "7", "8", "c"},
		{"Who was the first person to walk on the moon?", "Buzz Aldrin", "Yuri Gagarin", "Michael Collins", "Neil Armstrong", "d"},
		{"What is the largest desert in the world?", "Sahara Desert", "Arabian Desert", "Gobi Desert", "Antarctic Polar Desert", "d"},
		{"Which is the longest river in the world?", "Amazon River", "Nile River", "Yangtze River", "Mississippi River", "b"},
		{"What does 'CPU' stand for?", "Central Process Unit", "Computer Personal Unit", "Central Processing Unit", "Computer Primary Unit", "c"},
		{"Who wrote the play 'Romeo and Juliet'?", "Charles Dickens", "William Shakespeare", "George Orwell", "Jane Austen", "b"},
		{"What is the chemical symbol for gold?", "Ag", "Au", "Pb", "Fe", "b"},
		{"Which is the smallest planet in our solar system?", "Venus", "Mars", "Mercury", "Uranus", "c"},
		{"What is the capital of Canada?", "Toronto", "Vancouver", "Montreal", "Ottawa", "d"},
		{"How many bones are in the adult human body?", "206", "208", "210", "212", "a"},
		{"Which artist is known for the 'Starry Night' painting?", "Pablo Picasso", "Claude Monet", "Salvador Dalí", "Vincent van Gogh", "d"},
		{"What is the primary language spoken in Brazil?", "Spanish", "Portuguese", "Brazilian", "English", "b"},
		{"Who invented the telephone?", "Thomas Edison", "Nikola Tesla", "Alexander Graham Bell", "Guglielmo Marconi", "c"},
		{"What is the capital of Egypt?", "Alexandria", "Giza", "Cairo", "Luxor", "c"},
		{"Which gas do plants absorb from the atmosphere?", "Oxygen", "Nitrogen", "Carbon Dioxide", "Hydrogen", "c"},
		{"What is the main component of the sun?", "Liquid lava", "Rock", "Hydrogen and Helium", "Oxygen", "c"},
		{"Who was the first female Prime Minister of the United Kingdom?", "Theresa May", "Margaret Thatcher", "Angela Merkel", "Indira Gandhi", "b"},
		{"What is the largest country by land area?", "Canada", "China", "USA", "Russia", "d"},
		{"Which of these is a primary color?", "Green", "Orange", "Blue", "Purple", "c"},
		{"What is the boiling point of water at sea level?", "90°C", "100°C", "110°C", "120°C", "b"},
		{"Who is the author of the Harry Potter series?", "J.R.R. Tolkien", "George R.R. Martin", "Suzanne Collins", "J.K. Rowling", "d"},
		{"What is the capital of Italy?", "Milan", "Naples", "Rome", "Venice", "c"},
		{"Which ocean is the Bermuda Triangle located in?", "Atlantic", "Pacific", "Indian", "Arctic", "a"},
		{"What is the square root of 64?", "6", "7", "8", "9", "c"},
		{"What is the national sport of Canada?", "Hockey and Lacrosse", "Baseball", "Basketball", "Soccer", "a"},
		{"Which composer was deaf for the last part of his life?", "Mozart", "Bach", "Beethoven", "Chopin", "c"},
		{"What is the capital of Spain?", "Barcelona", "Seville", "Lisbon", "Madrid", "d"},
		{"What is the largest mammal in the world?", "Elephant", "Blue Whale", "Giraffe", "Great White Shark", "b"},
		{"Which country gifted the Statue of Liberty to the USA?", "Germany", "United Kingdom", "France", "Italy", "c"},
		{"What is the chemical formula for table salt?", "H2O", "C6H12O6", "NaCl", "CO2", "c"},
		{"Who discovered gravity?", "Albert Einstein", "Galileo Galilei", "Isaac Newton", "Nikola Tesla", "c"},
		{"What is the capital of Russia?", "Saint Petersburg", "Kazan", "Novosibirsk", "Moscow", "d"},
		{"In what year did World War II end?", "1943", "1944", "1945", "1946", "c"},
		{"What is the most spoken language in the world?", "English", "Spanish", "Mandarin Chinese", "Hindi", "c"},
		{"Who is known as the 'Father of Computers'?", "Alan Turing", "Charles Babbage", "Tim Berners-Lee", "Steve Jobs", "b"},
		{"What is the capital of India?", "Mumbai", "Kolkata", "Chennai", "New Delhi", "d"},
		{"What is the smallest continent by land area?", "South America", "Europe", "Antarctica", "Australia", "d"},
		{"Which famous scientist developed the theory of relativity?", "Isaac Newton", "Galileo Galilei", "Albert Einstein", "Stephen Hawking", "c"},
		{"What is the capital of China?", "Shanghai", "Hong Kong", "Beijing", "Tianjin", "c"},
		{"How many players are on a standard soccer team on the field?", "9", "10", "11", "12", "c"},
		{"What is the name of the galaxy we live in?", "Andromeda", "Triangulum", "Whirlpool", "Milky Way", "d"},
		{"Which is the largest bone in the human body?", "Tibia", "Humerus", "Femur", "Fibula", "c"},
		{"What is the capital of Germany?", "Munich", "Hamburg", "Frankfurt", "Berlin", "d"},
		{"Who was the first President of the United States?", "Thomas Jefferson", "Abraham Lincoln", "George Washington", "John Adams", "c"},
		{"What is the freezing point of water in Celsius?", "-10°C", "0°C", "10°C", "32°C", "b"},
		{"Which country is home to the kangaroo?", "New Zealand", "South Africa", "Australia", "Indonesia", "c"},
		{"What is the capital of Brazil?", "Rio de Janeiro", "São Paulo", "Salvador", "Brasília", "d"},
		{"What is the study of earthquakes called?", "Seismology", "Geology", "Volcanology", "Meteorology", "a"},
		{"Who painted the ceiling of the Sistine Chapel?", "Raphael", "Donatello", "Leonardo da Vinci", "Michelangelo", "d"},
		{"What is the main gas found in the air we breathe?", "Oxygen", "Carbon Dioxide", "Nitrogen", "Argon", "c"},
		{"What is the capital of Argentina?", "Santiago", "Lima", "Bogotá", "Buenos Aires", "d"},
		{"Which instrument is used to measure atmospheric pressure?", "Thermometer", "Barometer", "Hygrometer", "Anemometer", "b"},
		{"What is the largest planet in our solar system?", "Saturn", "Jupiter", "Neptune", "Uranus", "b"},
		{"Who wrote 'Pride and Prejudice'?", "Emily Brontë", "Charlotte Brontë", "Jane Austen", "Mary Shelley", "c"},
		{"What is the capital of South Korea?", "Busan", "Incheon", "Seoul", "Daegu", "c"},
		{"What is the process by which plants make their own food called?", "Respiration", "Transpiration", "Photosynthesis", "Germination", "c"},
		{"What is the world's longest man-made structure?", "The Great Wall of China", "The Hoover Dam", "The Panama Canal", "The Burj Khalifa", "a"},
		{"What is the capital of Mexico?", "Guadalajara", "Tijuana", "Cancún", "Mexico City", "d"},
		{"Which is the most populous country in the world?", "India", "United States", "Indonesia", "China", "a"},
		{"What is the chemical symbol for iron?", "I", "Ir", "Fe", "In", "c"},
		{"Who invented the light bulb?", "Nikola Tesla", "Benjamin Franklin", "Thomas Edison", "Alexander Graham Bell", "c"},
		{"What is the capital of South Africa?", "Cape Town", "Johannesburg", "Durban", "Pretoria", "d"},
		{"What type of animal is a 'canine'?", "Cat", "Dog", "Bird", "Fish", "b"},
		{"What is the main currency of the United Kingdom?", "Euro", "Dollar", "Pound Sterling", "Franc", "c"},
		{"Who was the ancient Greek god of the sea?", "Zeus", "Hades", "Poseidon", "Apollo", "c"},
		{"What is the capital of Thailand?", "Phuket", "Chiang Mai", "Pattaya", "Bangkok", "d"},
		{"In which country would you find the pyramids of Giza?", "Sudan", "Libya", "Egypt", "Jordan", "c"},
		{"What is the name of the world's largest rainforest?", "Congo Rainforest", "Daintree Rainforest", "Valdivian Rainforest", "The Amazon", "d"},
		{"What is the capital of Turkey?", "Istanbul", "Ankara", "Izmir", "Bursa", "b"},
		{"Which of the following is a reptile?", "Frog", "Snake", "Fish", "Bird", "b"},
		{"What is the chemical symbol for silver?", "Si", "Sv", "Ag", "Au", "c"},
		{"Who wrote 'The Great Gatsby'?", "Ernest Hemingway", "William Faulkner", "F. Scott Fitzgerald", "John Steinbeck", "c"},
		{"What is the capital of Greece?", "Thessaloniki", "Patras", "Heraklion", "Athens", "d"},
		{"What is the most common blood type in humans?", "A+", "B-", "O+", "AB+", "c"},
		{"Which country is famous for its tulips and windmills?", "Belgium", "Denmark", "Netherlands", "Switzerland", "c"},
		{"What is the capital of Sweden?", "Oslo", "Copenhagen", "Helsinki", "Stockholm", "d"},
		{"What is the largest bird in the world?", "Emu", "Ostrich", "Albatross", "Condor", "b"},
		{"Who is the main character in the 'Lord of the Rings' trilogy?", "Gandalf", "Aragorn", "Frodo Baggins", "Legolas", "c"},
		{"What is the capital of Norway?", "Bergen", "Trondheim", "Stavanger", "Oslo", "d"},
		{"What is the name of the force that opposes motion?", "Gravity", "Inertia", "Friction", "Momentum", "c"},
	}
}
