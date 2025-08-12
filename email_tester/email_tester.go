package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	baseURL      = "http://localhost:8070"
	numRequests  = 50
	maxSessionID = 30
)

func makeEmailRequest(wg *sync.WaitGroup, successCount *int32, failureCount *int32, requestNum int) {
	defer wg.Done()

	sessionID := rand.Intn(maxSessionID) + 1
	url := fmt.Sprintf("%s/sessions/%d/email-report", baseURL, sessionID)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		log.Printf("ERROR: Request #%d for session %d failed: %v", requestNum, sessionID, err)
		atomic.AddInt32(failureCount, 1)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		atomic.AddInt32(successCount, 1)
		log.Printf("SUCCESS: Email request #%d for session %d accepted (Status: %d)", requestNum, sessionID, resp.StatusCode)
	} else {
		atomic.AddInt32(failureCount, 1)
		log.Printf("FAILURE: Request #%d for session %d failed (Status: %d)", requestNum, sessionID, resp.StatusCode)
	}
}

func main() {
	// --- FIX STARTS HERE ---
	// Add a short delay to give the API server time to initialize the database.
	log.Println("Waiting 3 seconds for the API server to start and seed the database...")
	time.Sleep(3 * time.Second)
	// --- FIX ENDS HERE ---

	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	var successCount, failureCount int32

	log.Printf("Starting email load test with %d concurrent requests...", numRequests)
	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeEmailRequest(&wg, &successCount, &failureCount, i)
	}

	log.Println("All requests sent. Waiting for responses...")
	wg.Wait()

	duration := time.Since(startTime)
	log.Println("----------- Email Load Test Complete -----------")
	log.Printf("Total time taken: %s", duration)
	log.Printf("Successful requests: %d", successCount)
	log.Printf("Failed requests: %d", failureCount)
	log.Printf("Check your MailHog UI at http://localhost:8025 to see the emails.")
	log.Println("----------------------------------------------")
}
