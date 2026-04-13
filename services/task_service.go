package services

import (
	"context"
	"fmt"
	"task-tracker-api/models"
	"task-tracker-api/repository"
	"task-tracker-api/utils"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskServices struct {
	repo *repository.TaskRepository
}

var validate = validator.New()

func NewTaskServices(repo *repository.TaskRepository) *TaskServices {
	return &TaskServices{repo: repo}
}

func (s *TaskServices) CreateTask(ctx context.Context, task *models.Task) error {
	// validate
	if err := validate.Struct(task); err != nil {
		return err
	}

	// validate priority
	validPriorities := map[string]bool{"low": true, "medium": true, "high": true}
	if task.Priority != "" && !validPriorities[task.Priority] {
		return utils.BadRequest("priority must be low, medium or high")
	}

	// set default priority
	if task.Priority == "" {
		task.Priority = "low"
	}

	return s.repo.CreateTask(ctx, task)
}

func (s *TaskServices) GetTasks(ctx context.Context, filter *models.TaskFilter) (*models.PaginatedTasks, error) {
	// validate board id
	if filter.BoardID == "" {
		return nil, utils.BadRequest("board_id is required")
	}

	// validate priority
	if filter.Priority != "" {
		validPriorities := map[string]bool{"low": true, "medium": true, "high": true}
		if !validPriorities[filter.Priority] {
			return nil, utils.BadRequest("priority must be low, medium or high")
		}
	}

	// set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	// cap limit at 100
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	tasks, err := s.repo.GetTasks(ctx, filter)
	if err != nil {
		return nil, utils.InternalServerError(err.Error())
	}

	return tasks, nil
}

func (s *TaskServices) GetTask(ctx context.Context, id string) (*models.TaskResponse, error) {
	task, err := s.repo.GetTask(ctx, id)
	if err != nil {
		if err.Error() == "task not found" {
			return nil, utils.NotFound("task not found")
		}
		return nil, utils.InternalServerError("error fetching task")
	}

	return task, nil
}

// services/task_service.go
func (s *TaskServices) UpdateTask(ctx context.Context, id string, req *models.UpdateTask) error {
	if id == "" {
		return utils.BadRequest("task id is required")
	}

	if req.Title != nil && *req.Title == "" {
		return utils.BadRequest("title cannot be empty")
	}

	if req.Priority != nil {
		validPriorities := map[string]bool{"low": true, "medium": true, "high": true}
		if !validPriorities[*req.Priority] {
			return utils.BadRequest("priority must be low, medium or high")
		}
	}

	if err := s.repo.UpdateTask(ctx, id, req); err != nil {
		if err.Error() == "task not found" {
			return utils.NotFound("task not found")
		}
		return utils.InternalServerError("error updating task")
	}

	return nil
}

func (s *TaskServices) Column(ctx context.Context, id string, req *models.Column) error {
	return s.repo.Column(ctx, id, req)
}

func (s *TaskServices) CreateColumn(ctx context.Context, req *models.CreateColumn) (*models.Column, error) {
	if req.BoardID == "" {
		return nil, utils.BadRequest("board_id is required")
	}

	validNames := map[string]bool{
		"todo":        true,
		"in progress": true,
		"completed":   true,
	}
	if !validNames[req.Name] {
		return nil, utils.BadRequest("name must be todo, in progress or completed")
	}

	// convert string board_id to ObjectID
	boardObjectID, err := bson.ObjectIDFromHex(req.BoardID)
	if err != nil {
		return nil, utils.BadRequest("invalid board_id")
	}

	col := &models.Column{
		BoardID: boardObjectID,
		Name:    req.Name,
	}

	if err := s.repo.CreateColumn(ctx, col); err != nil {
		return nil, utils.InternalServerError("error creating column")
	}

	return col, nil
}

func (s *TaskServices) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	return s.repo.DeleteTask(ctx, id)

}

func (s *TaskServices) UpdateColumn(ctx context.Context, id string, req *models.UpdateColumn) error {
	if id == "" {
		return utils.BadRequest("column id is required")
	}

	if req.Name == nil || *req.Name == "" {
		return utils.BadRequest("column name is required")
	}

	// validate column name
	validNames := map[string]bool{
		"todo":        true,
		"in progress": true,
		"completed":   true,
	}
	if !validNames[*req.Name] {
		return utils.BadRequest("column name must be todo, in progress or completed")
	}

	if err := s.repo.UpdateColumn(ctx, id, req); err != nil {
		if err.Error() == "column not found" {
			return utils.NotFound("column not found")
		}
		return utils.InternalServerError("error updating column")
	}

	return nil
}
