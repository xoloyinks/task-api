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

func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) error {
	var req models.TeamMember

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.BadRequest(err.Error())
	}

	// get teamID from URL path instead of request body
	teamID := r.PathValue("id")
	email := req.Email
	role := req.Role

	if err := h.service.AddMember(r.Context(), teamID, email, role); err != nil {
		return utils.InternalServerError(err.Error())
	}

	return utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "member added to team successfully",
	})
}

func (h *TeamHandler) GetAllTeams(w http.ResponseWriter, r *http.Request) error {
	teams, err := h.service.GetAllTeams(r.Context())
	if err != nil {
		return utils.InternalServerError(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)

	return nil
}
