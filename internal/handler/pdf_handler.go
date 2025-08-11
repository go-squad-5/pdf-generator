package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	idStr, ok := vars["id"]
	if !ok {
		models.RespondWithError(w, http.StatusBadRequest, "Session ID is missing")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		models.RespondWithError(w, http.StatusBadRequest, "Invalid Session ID format")
		return
	}

	log.Printf("Received request to generate report for session ID: %d", id)

	pdfBytes, err := h.service.GenerateQuizReport(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("session with ID %d not found", id) {
			models.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			models.RespondWithError(w, http.StatusInternalServerError, "Failed to generate PDF report")
		}
		log.Printf("Error generating report for session ID %d: %v", id, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=quiz_report_session_"+idStr+".pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}
