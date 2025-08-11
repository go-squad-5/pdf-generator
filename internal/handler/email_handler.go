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

type EmailHandler struct {
	service *service.EmailService
}

func NewEmailHandler(s *service.EmailService) *EmailHandler {
	return &EmailHandler{service: s}
}

func (h *EmailHandler) SendReportHandler(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Received request to email detailed report for session ID: %d", id)

	err = h.service.SendQuizReportByEmail(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("session with ID %d not found", id) {
			models.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			models.RespondWithError(w, http.StatusInternalServerError, "Failed to process email request")
		}
		log.Printf("Error processing email for session ID %d: %v", id, err)
		return
	}

	response := map[string]string{"message": fmt.Sprintf("Detailed quiz report for session %d is being sent in parts.", id)}
	models.RespondWithJSON(w, http.StatusAccepted, response)
}
