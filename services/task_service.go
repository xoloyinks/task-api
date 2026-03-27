package services

import (
	"context"
	"fmt"
	"task-tracker-api/models"
	"task-tracker-api/repository"
	"task-tracker-api/utils"

	"github.com/go-playground/validator/v10"
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

func (s *TaskServices) GetTasks(ctx context.Context, boardID string) ([]models.TaskResponse, error) {
	if boardID == "" {
		return nil, utils.BadRequest("board id is required")
	}

	tasks, err := s.repo.GetTasks(ctx, boardID)
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
