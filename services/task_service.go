package services

import (
	"context"
	"fmt"
	"task-tracker-api/models"
	"task-tracker-api/repository"

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
	if task.Title == "" {
		return fmt.Errorf("Title required for task")
	}
	if len(task.Title) > 100 {
		return fmt.Errorf("title cannot exceed 100 characters")
	}
	task.Completed = false

	return s.repo.CreateTask(ctx, task)
}

func (s *TaskServices) GetAllTasks(ctx context.Context) (error, []models.Task) {
	return s.repo.GetAllTasks(ctx)
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

func (s *TaskServices) CompleteTask(ctx context.Context, id string, req *models.CompleteTask) error {
	return s.repo.CompleteTask(ctx, id, req)
}

func (s *TaskServices) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	return s.repo.DeleteTask(ctx, id)

}
