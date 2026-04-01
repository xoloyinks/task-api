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

// handlers/team_handler.go
func (h *TeamHandler) GetTeams(w http.ResponseWriter, r *http.Request) error {
	claims := r.Context().Value(utils.ClaimsKey).(*utils.Claims)

	teams, err := h.service.GetTeams(r.Context(), claims.UserID)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, teams)
}

// handlers/team_handler.go
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	team, err := h.service.GetTeam(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, team)
}

// handlers/team_handler.go
func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	var req models.UpdateTeam
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.BadRequest("invalid request body")
	}

	if err := h.service.UpdateTeam(r.Context(), id, &req); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "team updated successfully",
	})
}

// handlers/team_handler.go
func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	if err := h.service.DeleteTeam(r.Context(), id); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "team and all its data deleted successfully",
	})
}

// handlers/team_handler.go
func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) error {
	teamID := r.PathValue("id")
	memberID := r.PathValue("memberID")

	if err := h.service.RemoveMember(r.Context(), teamID, memberID); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "member removed successfully",
	})
}
