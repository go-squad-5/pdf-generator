package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mockEmailService struct {
	err error
}

func (m *mockEmailService) SendQuizReportByEmail(sessionID string) error {
	return m.err
}

func TestEmailHandler_Success(t *testing.T) {
	mockService := &mockEmailService{
		err: nil,
	}
	handler := NewEmailHandler(mockService)

	req, _ := http.NewRequest("POST", "/sessions/valid-ssid/email-report", nil)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()

	router.HandleFunc("/sessions/{id}/email-report", handler.SendReportHandler)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code, "Status code should be 202 Accepted")
	assert.Contains(t, rr.Body.String(), "is being sent", "The success message is incorrect")
}

func TestEmailHandler_SessionIdNotFound(t *testing.T) {
	sessionID := "session-id-invalid"
	mockService := &mockEmailService{
		err: fmt.Errorf("session with ID %s not found", sessionID),
	}

	handler := NewEmailHandler(mockService)

	req, _ := http.NewRequest("POST", "/sessions/session-id-invalid/email-report", nil)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()

	router.HandleFunc("/sessions/{id}/email-report", handler.SendReportHandler)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "session with ID session-id-invalid not found")

}

func TestEmailHandler_BadRequest(t *testing.T) {
	mockService := &mockEmailService{
		err: errors.New("internal server error"),
	}

	handler := NewEmailHandler(mockService)

	req, _ := http.NewRequest("POST", "/sessions/some-session-id/email-report", nil)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()

	router.HandleFunc("/sessions/{id}/email-report", handler.SendReportHandler)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestEmailHandler_MissingSessionId(t *testing.T) {
	handler := NewEmailHandler(nil)

	req, _ := http.NewRequest("POST", "/sessions//email-report", nil)

	rr := httptest.NewRecorder()

	handler.SendReportHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be 400 Bad Request")
	assert.Contains(t, rr.Body.String(), "Session ID is missing", "Response body should contain the correct error message")
}
