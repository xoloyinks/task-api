package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"task-tracker-api/models"
	"task-tracker-api/services"
	"task-tracker-api/sse"
	"task-tracker-api/utils"
)

type TaskHandler struct {
	service *services.TaskServices
	hub     *sse.Hub
}

func NewTaskHandler(service *services.TaskServices, hub *sse.Hub) *TaskHandler {
	return &TaskHandler{service: service, hub: hub}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) error {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return utils.BadRequest("invalid request body")
	}

	claims := r.Context().Value(utils.ClaimsKey).(*utils.Claims)
	task.DestinationID = claims.UserID

	if err := h.service.CreateTask(r.Context(), &task); err != nil {
		return err
	}

	// broadcast to all clients in the team
	taskJSON, _ := json.Marshal(task)
	h.hub.Broadcast(task.BoardID.Hex(), "task_created", string(taskJSON))

	return utils.WriteJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	var req models.UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.BadRequest("invalid request body")
	}

	if err := h.service.UpdateTask(r.Context(), id, &req); err != nil {
		return err
	}

	// broadcast update
	reqJSON, _ := json.Marshal(req)
	h.hub.Broadcast(id, "task_updated", string(reqJSON))

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "task updated successfully",
	})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		return err
	}

	// broadcast delete
	h.hub.Broadcast(id, "task_deleted", fmt.Sprintf(`{"id": "%s"}`, id))

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "task deleted successfully",
	})
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, task)
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

func (h *TaskHandler) UpdateColumn(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	var req models.UpdateColumn
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.BadRequest("invalid request body")
	}

	if err := h.service.UpdateColumn(r.Context(), id, &req); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "column updated successfully",
	})
}
