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

func (s *TeamServices) GetAllTeams(ctx context.Context) ([]models.Team, error) {

	return s.repo.GetAllTeams(ctx)

}
