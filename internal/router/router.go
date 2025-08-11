package router

import (
	"github.com/go-squad-5/pdf-generator/internal/handler"
	"github.com/gorilla/mux"
)

func NewRouter(pdfHandler *handler.PDFHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users/{id}/report", pdfHandler.GenerateReportHandler).Methods("GET")

	return r
}
