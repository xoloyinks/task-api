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

func (s *TaskServices) GetAllTasks(ctx context.Context, userId string) (error, []models.Task) {
	return s.repo.GetAllTasks(ctx, userId)
}

func (s *TaskServices) GetTask(ctx context.Context, id string) (*models.Task, error) {

	if id == "" {
		return nil, fmt.Errorf("invalid task id")
	}

	return s.repo.GetTask(ctx, id)

}

func (s *TaskServices) UpdateTask(ctx context.Context, id string, req *models.UpdateTask) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}

	return s.repo.UpdateTask(ctx, id, req)
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
