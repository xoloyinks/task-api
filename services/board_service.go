package services

import (
	"context"
	"task-tracker-api/models"
	"task-tracker-api/repository"
	"task-tracker-api/utils"
)

type BoardServices struct {
	repo *repository.BoardRepository
}

func NewBoardServices(repo *repository.BoardRepository) *BoardServices {
	return &BoardServices{repo: repo}
}

func (s *BoardServices) CreateBoard(ctx context.Context, req *models.Board) error {
	if req.DestinationID == "" {
		return utils.BadRequest("Destination id required")
	}

	if req.Name == "" {
		return utils.BadRequest("Board name is required")
	}

	return s.repo.CreateBoard(ctx, req)
}

// services/board_service.go
func (s *BoardServices) GetBoard(ctx context.Context, id string) (*models.Board, error) {
	board, err := s.repo.GetBoard(ctx, id)
	if err != nil {
		if err.Error() == "board not found" {
			return nil, utils.NotFound("board not found")
		}
		return nil, utils.InternalServerError("error fetching board")
	}
	return board, nil
}

// services/board_service.go
func (s *BoardServices) GetAllBoards(ctx context.Context) ([]models.Board, error) {
	boards, err := s.repo.GetAllBoards(ctx)
	if err != nil {
		return nil, utils.InternalServerError("error fetching boards")
	}
	return boards, nil
}
