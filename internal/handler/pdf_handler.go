package handler

import (
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
		models.RespondWithError(w, http.StatusBadRequest, "User ID is missing")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		models.RespondWithError(w, http.StatusBadRequest, "Invalid User ID format")
		return
	}

	log.Printf("Received request to generate report for user ID: %d", id)

	pdfBytes, err := h.service.GenerateUserReport(id)
	if err != nil {
		if err.Error() == "user with ID "+idStr+" not found" {
			models.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			models.RespondWithError(w, http.StatusInternalServerError, "Failed to generate PDF report")
		}
		log.Printf("Error generating report for user ID %d: %v", id, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=user_report_"+idStr+".pdf")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(pdfBytes)
	if err != nil {
		log.Printf("Failed to write PDF response for user ID %d: %v", id, err)
	}
	log.Printf("Successfully sent PDF report for user ID: %d", id)
}
