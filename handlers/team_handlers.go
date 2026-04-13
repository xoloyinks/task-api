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

// CreateTeam godoc
// @Summary      Create team
// @Description  Create a new team
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        team body     models.Team  true  "Team payload"
// @Success      201  {object} map[string]string
// @Failure      400  {object} utils.AppError
// @Failure      500  {object} utils.AppError
// @Router       /teams [post]
func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) error {
	var req models.Team

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return utils.BadRequest(err.Error())

	}

	teamId, err := h.service.CreateTeam(r.Context(), &req)

	if err != nil {
		return utils.InternalServerError(err.Error())

	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Team created",
		"teamId":  teamId,
	})

	return nil
}

// AddMember godoc
// @Summary      Add member to team
// @Description  Add a new member to a team using email and role
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        id     path      string              true  "Team ID"
// @Param        member body      models.TeamMember   true  "Member payload"
// @Success      201    {object}  map[string]string
// @Failure      400    {object}  utils.AppError
// @Failure      500    {object}  utils.AppError
// @Router       /teams/{id}/members [post]
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

// GetTeams godoc
// @Summary      Get user teams
// @Description  Retrieve all teams for the authenticated user
// @Tags         teams
// @Produce      json
// @Success      200  {array}   models.Team
// @Failure      401  {object}  utils.AppError
// @Failure      500  {object}  utils.AppError
// @Router       /teams [get]
// handlers/team_handler.go
func (h *TeamHandler) GetTeams(w http.ResponseWriter, r *http.Request) error {
	claims := r.Context().Value(utils.ClaimsKey).(*utils.Claims)

	teams, err := h.service.GetTeams(r.Context(), claims.Email)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, teams)
}

// GetTeam godoc
// @Summary      Get team
// @Description  Retrieve a single team by ID
// @Tags         teams
// @Produce      json
// @Param        id   path      string  true  "Team ID"
// @Success      200  {object}  models.Team
// @Failure      404  {object}  utils.AppError
// @Failure      500  {object}  utils.AppError
// @Router       /teams/{id} [get]
// handlers/team_handler.go
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	team, err := h.service.GetTeam(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, team)
}

// UpdateTeam godoc
// @Summary      Update team
// @Description  Update a team by ID
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "Team ID"
// @Param        team  body      models.UpdateTeam   true  "Updated team payload"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  utils.AppError
// @Failure      404   {object}  utils.AppError
// @Failure      500   {object}  utils.AppError
// @Router       /teams/{id} [put]
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

// DeleteTeam godoc
// @Summary      Delete team
// @Description  Delete a team and all associated data
// @Tags         teams
// @Produce      json
// @Param        id   path      string  true  "Team ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  utils.AppError
// @Failure      500  {object}  utils.AppError
// @Router       /teams/{id} [delete]
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

// RemoveMember godoc
// @Summary      Remove team member
// @Description  Remove a member from a team
// @Tags         teams
// @Produce      json
// @Param        id        path      string  true  "Team ID"
// @Param        memberID  path      string  true  "Member ID"
// @Success      200       {object}  map[string]string
// @Failure      404       {object}  utils.AppError
// @Failure      500       {object}  utils.AppError
// @Router       /teams/{id}/members/{memberID} [delete]
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
