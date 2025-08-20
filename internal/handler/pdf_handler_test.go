package handler

import (
  "fmt"
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/gorilla/mux"
  "github.com/stretchr/testify/assert"
)

type mockPDFService struct {
  pdfBytes []byte
  err      error
}

func (m *mockPDFService) GenerateQuizReport(sessionID string) ([]byte, error) {
  return m.pdfBytes, m.err
}

func TestPDFHandler_GenerateReportHandler_Success(t *testing.T) {

  mockService := &mockPDFService{
    pdfBytes: []byte("pdf-content-idhr-aayega"),
    err:      nil,
  }

  handler := NewPDFHandler(mockService)

  req, _ := http.NewRequest("GET", "/sessions/sid-123/report", nil)

  rr := httptest.NewRecorder()

  router := mux.NewRouter()
  router.HandleFunc("/sessions/{id}/report", handler.GenerateReportHandler)

  router.ServeHTTP(rr, req)

  assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200 OK")
  assert.Equal(t, "application/pdf", rr.Header().Get("Content-Type"), "Content-Type header should be application/pdf")
  assert.Equal(t, "attachment; filename=quiz_report_session_sid-123.pdf", rr.Header().Get("Content-Disposition"), "Content-Disposition header is incorrect")
  assert.Equal(t, "pdf-content-idhr-aayega", rr.Body.String(), "The response body should contain the PDF content")
}

func TestPDFHandler_MissingID(t *testing.T) {
  handler := NewPDFHandler(nil)

  req, _ := http.NewRequest("GET", "/sessions//report", nil)

  rr := httptest.NewRecorder()

  handler.GenerateReportHandler(rr, req)

  assert.Equal(t, http.StatusBadRequest, rr.Code, "The status code should be 400 for a bad request")
  assert.Contains(t, rr.Body.String(), "Session ID is missing", "The response body should contain the missing ID error")
}

func TestPDFHandler_SessionIDNotFound(t *testing.T) {
  mockService := &mockPDFService{
    pdfBytes: nil,
    err:      fmt.Errorf("session with ID %s not found", "some-id"),
  }

  handler := NewPDFHandler(mockService)

  req, _ := http.NewRequest("GET", "/sessions/some-id/report", nil)

  rr := httptest.NewRecorder()

  router := mux.NewRouter()
  router.HandleFunc("/sessions/{id}/report", handler.GenerateReportHandler)

  router.ServeHTTP(rr, req)

  assert.Equal(t, http.StatusNotFound, rr.Code, "The status code should be 404 for a not found error")
  assert.Contains(t, rr.Body.String(), "session with ID some-id not found", "The response body should contain the not found error message")
}
