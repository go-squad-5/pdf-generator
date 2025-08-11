package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	baseURL      = "http://localhost:8080"
	numRequests  = 100
	maxSessionID = 30
	reportsDir   = "./reports"
)

func makeRequest(wg *sync.WaitGroup, successCount *int32, failureCount *int32, requestNum int) {
	defer wg.Done()

	sessionID := rand.Intn(maxSessionID) + 1
	url := fmt.Sprintf("%s/sessions/%d/report", baseURL, sessionID)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("ERROR: Request #%d for session %d failed: %v", requestNum, sessionID, err)
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

		filePath := fmt.Sprintf("%s/report_session_%d_req_%d.pdf", reportsDir, sessionID, requestNum)
		err = ioutil.WriteFile(filePath, pdfBytes, 0644)
		if err != nil {
			log.Printf("ERROR: Failed to save PDF for request #%d: %v", requestNum, err)
			atomic.AddInt32(failureCount, 1)
			return
		}

		atomic.AddInt32(successCount, 1)
		log.Printf("SUCCESS: Saved PDF for request #%d (session %d) to %s", requestNum, sessionID, filePath)
	} else {
		atomic.AddInt32(failureCount, 1)
		log.Printf("FAILURE: Request #%d for session %d failed (Status: %d)", requestNum, sessionID, resp.StatusCode)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if _, err := os.Stat(reportsDir); os.IsNotExist(err) {
		log.Printf("Creating directory: %s", reportsDir)
		os.Mkdir(reportsDir, os.ModePerm)
	}

	var wg sync.WaitGroup
	var successCount, failureCount int32

	log.Printf("Starting load test with %d concurrent requests...", numRequests)
	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeRequest(&wg, &successCount, &failureCount, i)
	}

	log.Println("All requests sent. Waiting for responses...")
	wg.Wait()

	duration := time.Since(startTime)
	log.Println("----------- Load Test Complete -----------")
	log.Printf("Total time taken: %s", duration)
	log.Printf("Successful requests (PDFs saved): %d", successCount)
	log.Printf("Failed requests: %d", failureCount)
	log.Printf("Check the '%s' folder for the generated PDF files.", reportsDir)
	log.Println("----------------------------------------")
}
