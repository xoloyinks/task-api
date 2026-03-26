package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"task-tracker-api/models"
	"task-tracker-api/services"
	"task-tracker-api/utils"
)

type TaskHandler struct {
	service *services.TaskServices
}

func NewTaskHandler(service *services.TaskServices) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask godoc
// @Summary      Create a task
// @Description  Add a new task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        task body models.Task true "Task object"
// @Success      201 {object} models.Task
// @Failure      400 {object} utils.AppError
// @Security     BearerAuth
// @Router       /task [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) error {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return utils.BadRequest("invalid request body")
	}

	// get logged in user from context
	claims := r.Context().Value(utils.ClaimsKey).(*utils.Claims)
	task.DestinationID = claims.UserID

	if err := h.service.CreateTask(r.Context(), &task); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, task)
}

// GetAllTasks godoc
// @Summary      Get all tasks
// @Description  Get all tasks for the logged in user
// @Tags         tasks
// @Produce      json
// @Param        completed query bool false "Filter by completed status"
// @Success      200 {array} models.Task
// @Failure      500 {object} utils.AppError
// @Security     BearerAuth
// @Router       /tasks [get]
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) error {
	userId := r.PathValue("userId")
	err, tasks := h.service.GetAllTasks(r.Context(), userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

	return nil
}

// GetTask godoc
// @Summary      Get a task
// @Description  Get a single task by ID
// @Tags         tasks
// @Produce      json
// @Param        id path string true "Task ID"
// @Success      200 {object} models.Task
// @Failure      404 {object} utils.AppError
// @Security     BearerAuth
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	task, err := h.service.GetTask(r.Context(), id)

	if err != nil {
		if err.Error() == "invalid task id" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return err
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)

	return nil
}

func (h *TaskHandler) CheckId(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return err
	}

	msg := fmt.Sprintf("id value is %d", id)

	w.Write([]byte(msg))

	return nil

}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	var req models.UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return err
	}

	if err := h.service.UpdateTask(r.Context(), id, &req); err != nil {
		if err.Error() == "task not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return err
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "task updated successfully"})

	return nil
}

func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	var req *models.Column
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return err
	}

	if err := h.service.Column(r.Context(), id, req); err != nil {
		if err.Error() == "task not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return err
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "task status updated successfully"})

	return nil
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		if err.Error() == "id cannot be empty" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "task deleted successfully"})

	return nil
}
