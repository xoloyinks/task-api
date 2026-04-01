package services

import (
	"context"
	"fmt"
	"task-tracker-api/models"
	"task-tracker-api/repository"
	"task-tracker-api/utils"
)

type TeamServices struct {
	repo     *repository.TeamRepository
	authRepo *repository.AuthRepository
}

func NewTeamServices(repo *repository.TeamRepository) *TeamServices {
	return &TeamServices{repo: repo}
}

func (s *TeamServices) CreateTeam(ctx context.Context, req *models.Team) error {
	if req.CreatedBy == "" {
		return fmt.Errorf("Creator id is required")
	}

	existingTeam, _ := s.repo.GetTeamByName(ctx, req.Name)
	if existingTeam != nil {
		return utils.BadRequest("Team with name already exists")
	}

	return s.repo.CreateTeam(ctx, req)
}

func (s *TeamServices) AddMember(ctx context.Context, teamID string, email string, role string) error {
	if teamID == "" {
		return utils.BadRequest("team id is required")
	}

	if email == "" {
		return utils.BadRequest("email is required")
	}

	// validate role
	validRoles := map[string]bool{"admin": true, "member": true, "viewer": true}
	if !validRoles[role] {
		return utils.BadRequest("role must be admin, member or viewer")
	}

	// find user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return utils.InternalServerError("error finding user")
	}

	// check if user exists
	if user == nil {
		return utils.NotFound("user not found")
	}

	// check if user is already a member
	for _, id := range user.TeamID {
		if id == teamID {
			return utils.BadRequest("user is already a member of this team")
		}
	}

	if err := s.repo.AddMember(ctx, teamID, email, role); err != nil {
		return err
	}

	return nil
}

// services/team_service.go
func (s *TeamServices) GetTeams(ctx context.Context, userID string) ([]models.Team, error) {
	teams, err := s.repo.GetTeams(ctx, userID)
	if err != nil {
		return nil, utils.InternalServerError("error fetching teams")
	}

	return teams, nil
}

// services/team_service.go
func (s *TeamServices) GetTeam(ctx context.Context, id string) (*models.TeamResponse, error) {
	team, err := s.repo.GetTeam(ctx, id)
	if err != nil {
		if err.Error() == "team not found" {
			return nil, utils.NotFound("team not found")
		}
		return nil, utils.InternalServerError("error fetching team")
	}

	return team, nil
}

// services/team_service.go
func (s *TeamServices) UpdateTeam(ctx context.Context, id string, req *models.UpdateTeam) error {
	if id == "" {
		return utils.BadRequest("team id is required")
	}

	if req.Name == nil && req.Description == nil {
		return utils.BadRequest("at least one field is required")
	}

	if req.Name != nil && *req.Name == "" {
		return utils.BadRequest("name cannot be empty")
	}

	if err := s.repo.UpdateTeam(ctx, id, req); err != nil {
		if err.Error() == "team not found" {
			return utils.NotFound("team not found")
		}
		if err.Error() == "no fields to update" {
			return utils.BadRequest("at least one field is required")
		}
		return utils.InternalServerError("error updating team")
	}

	return nil
}

// services/team_service.go
func (s *TeamServices) DeleteTeam(ctx context.Context, id string) error {
	if id == "" {
		return utils.BadRequest("team id is required")
	}

	if err := s.repo.DeleteTeam(ctx, id); err != nil {
		if err.Error() == "team not found" {
			return utils.NotFound("team not found")
		}
		return utils.InternalServerError("error deleting team")
	}

	return nil
}

// services/team_service.go
func (s *TeamServices) RemoveMember(ctx context.Context, teamID string, memberID string) error {
	if teamID == "" {
		return utils.BadRequest("team id is required")
	}

	if memberID == "" {
		return utils.BadRequest("member id is required")
	}

	if err := s.repo.RemoveMember(ctx, teamID, memberID); err != nil {
		if err.Error() == "member not found" {
			return utils.NotFound("member not found")
		}
		return utils.InternalServerError("error removing member")
	}

	return nil
}
