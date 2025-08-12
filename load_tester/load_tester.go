package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	baseURL     = "http://localhost:8070"
	numRequests = 50
	reportsDir  = "./reports"
)

func getSessionIDs() []string {
	dsn := "root:root@tcp(127.0.0.1:3333)/quizdb?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL to get session IDs: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT session_id FROM Session")
	if err != nil {
		log.Fatalf("Failed to query for session IDs: %v", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Printf("Warning: could not scan session ID: %v", err)
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func makeRequest(wg *sync.WaitGroup, successCount *int32, failureCount *int32, sessionID string, requestNum int) {
	defer wg.Done()

	url := fmt.Sprintf("%s/sessions/%s/report", baseURL, sessionID)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("ERROR: Request #%d for session %s failed: %v", requestNum, sessionID, err)
		atomic.AddInt32(failureCount, 1)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		pdfBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ERROR: Failed to read response body for request #%d: %v", requestNum, err)
			atomic.AddInt32(failureCount, 1)
			return
		}

		filePath := fmt.Sprintf("%s/report_session_%s_req_%d.pdf", reportsDir, sessionID, requestNum)
		err = ioutil.WriteFile(filePath, pdfBytes, 0644)
		if err != nil {
			log.Printf("ERROR: Failed to save PDF for request #%d: %v", requestNum, err)
			atomic.AddInt32(failureCount, 1)
			return
		}

		atomic.AddInt32(successCount, 1)
		log.Printf("SUCCESS: Saved PDF for request #%d (session %s)", requestNum, sessionID)
	} else {
		atomic.AddInt32(failureCount, 1)
		log.Printf("FAILURE: Request #%d for session %s failed (Status: %d)", requestNum, sessionID, resp.StatusCode)
	}
}

func main() {
	log.Println("Waiting 3 seconds for the API server to be ready...")
	time.Sleep(3 * time.Second)

	log.Println("Fetching session IDs from the database...")
	sessionIDs := getSessionIDs()
	if len(sessionIDs) == 0 {
		log.Fatal("No session IDs found in the database. Please run the seeder first.")
	}
	log.Printf("Found %d sessions to test with.", len(sessionIDs))

	if _, err := os.Stat(reportsDir); os.IsNotExist(err) {
		log.Printf("Creating directory: %s", reportsDir)
		os.Mkdir(reportsDir, os.ModePerm)
	}

	var wg sync.WaitGroup
	var successCount, failureCount int32

	log.Printf("Starting PDF load test with %d concurrent requests...", numRequests)
	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		randomSessionID := sessionIDs[rand.Intn(len(sessionIDs))]
		go makeRequest(&wg, &successCount, &failureCount, randomSessionID, i)
	}

	log.Println("All requests sent. Waiting for responses...")
	wg.Wait()

	duration := time.Since(startTime)
	log.Println("----------- PDF Load Test Complete -----------")
	log.Printf("Total time taken: %s", duration)
	log.Printf("Successful requests (PDFs saved): %d", successCount)
	log.Printf("Failed requests: %d", failureCount)
	log.Printf("Check the '%s' folder for the generated PDF files.", reportsDir)
	log.Println("--------------------------------------------")
}
