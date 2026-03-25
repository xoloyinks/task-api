package services

import (
	"context"
	"fmt"
	"task-tracker-api/models"
	"task-tracker-api/repository"
	"task-tracker-api/utils"
)

type TeamServices struct {
	repo *repository.TeamRepository
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
