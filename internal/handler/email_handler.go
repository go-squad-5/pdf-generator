package handler

import (
	"fmt"
	"net/http"

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
	// The session ID is now a string (UUID).
	id, ok := vars["id"]
	if !ok {
		models.RespondWithError(w, http.StatusBadRequest, "Session ID is missing")
		return
	}

	err := h.service.SendQuizReportByEmail(id)
	if err != nil {
		// Use the string ID in the error message.
		if err.Error() == fmt.Sprintf("session with ID %s not found", id) {
			models.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			models.RespondWithError(w, http.StatusInternalServerError, "Failed to process email request")
		}
		return
	}

	response := map[string]string{"message": fmt.Sprintf("Detailed quiz report for session %s is being sent.", id)}
	models.RespondWithJSON(w, http.StatusAccepted, response)
}
