package router

import (
	"github.com/go-squad-5/pdf-generator/internal/handler"
	"github.com/gorilla/mux"
)

func NewRouter(pdfHandler *handler.PDFHandler, emailHandler *handler.EmailHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/sessions/{id}/report", pdfHandler.GenerateReportHandler).Methods("GET")

	r.HandleFunc("/sessions/{id}/email-report", emailHandler.SendReportHandler).Methods("POST")

	return r
}
