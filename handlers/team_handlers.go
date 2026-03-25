package handlers

import (
	"encoding/json"
	"net/http"
	"task-tracker-api/models"
	"task-tracker-api/services"
	"task-tracker-api/utils"
)

type TeamHandler struct {
	service *services.TeamServices
}

func NewTeamHandler(service *services.TeamServices) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) error {
	var req models.Team

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return utils.BadRequest(err.Error())

	}

	if err := h.service.CreateTeam(r.Context(), &req); err != nil {
		return utils.InternalServerError(err.Error())

	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Team created",
	})

	return nil
}
