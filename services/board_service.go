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

// services/board_service.go
func (s *BoardServices) UpdateBoard(ctx context.Context, id string, req *models.UpdateBoard) error {
	if id == "" {
		return utils.BadRequest("board id is required")
	}

	if req.Name == nil {
		return utils.BadRequest("at least one field is required")
	}

	if err := s.repo.UpdateBoard(ctx, id, req); err != nil {
		if err.Error() == "board not found" {
			return utils.NotFound("board not found")
		}
		if err.Error() == "name cannot be empty" {
			return utils.BadRequest("name cannot be empty")
		}
		return utils.InternalServerError("error updating board")
	}

	return nil
}

// services/board_service.go
func (s *BoardServices) GetBoardsByDestination(ctx context.Context, destinationID string) ([]models.Board, error) {
	if destinationID == "" {
		return nil, utils.BadRequest("destination id is required")
	}

	boards, err := s.repo.GetBoardsByDestination(ctx, destinationID)
	if err != nil {
		return nil, utils.InternalServerError("error fetching boards")
	}

	return boards, nil
}

// services/board_service.go
func (s *BoardServices) DeleteBoard(ctx context.Context, id string) error {
	if id == "" {
		return utils.BadRequest("board id is required")
	}

	if err := s.repo.DeleteBoard(ctx, id); err != nil {
		if err.Error() == "board not found" {
			return utils.NotFound("board not found")
		}
		return utils.InternalServerError("error deleting board")
	}

	return nil
}
