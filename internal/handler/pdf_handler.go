package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-squad-5/pdf-generator/internal/models"
	"github.com/go-squad-5/pdf-generator/internal/service"
	"github.com/gorilla/mux"
)

type PDFHandler struct {
	service *service.PDFService
}

func NewPDFHandler(s *service.PDFService) *PDFHandler {
	return &PDFHandler{service: s}
}

func (h *PDFHandler) GenerateReportHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// The session ID is now a string (UUID), so we don't need to convert it to an integer.
	id, ok := vars["id"]
	if !ok {
		log.Println("ID is required, bad request.")
		models.RespondWithError(w, http.StatusBadRequest, "Session ID is missing")
		return
	}
	log.Printf("Generating PDF report for session ID: %s", id)

	pdfBytes, err := h.service.GenerateQuizReport(id)
	if err != nil {
		// Use the string ID in the error message.
		if err.Error() == fmt.Sprintf("session with ID %s not found", id) {
			log.Println("Error while fetching by session id")
			models.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			log.Println("Failed to generate PDF Report", err)
			models.RespondWithError(w, http.StatusInternalServerError, "Failed to generate PDF report")
		}
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=quiz_report_session_"+id+".pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}
