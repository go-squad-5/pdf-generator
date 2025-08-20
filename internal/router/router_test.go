package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-squad-5/pdf-generator/internal/handler"
	"github.com/stretchr/testify/assert"
)

type mockPDFService struct{}

func (m *mockPDFService) GenerateQuizReport(sessionID string) ([]byte, error) {
	return nil, nil
}

type mockEmailService struct{}

func (m *mockEmailService) SendQuizReportByEmail(sessionID string) error {
	return nil
}
func TestNewRouter(t *testing.T) {
	pdfHandler := handler.NewPDFHandler(&mockPDFService{})
	emailHandler := handler.NewEmailHandler(&mockEmailService{})
	router := NewRouter(pdfHandler, emailHandler)

	t.Run("TestNewRouterForPdfHandler", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/sessions/some-id/report", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})

	t.Run("TestNewRouterForEmailHandler", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/sessions/some-id/email-report", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.NotEqual(t, http.StatusNotFound, rr.Code)
	})
}
