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

// CreateTask godoc
// @Summary      Create a task
// @Description  Create a new task on a board. A default todo column is created automatically if one does not exist for the board.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task body models.Task true "Task object"
// @Success      201 {object} models.Task
// @Failure      400 {object} utils.AppError
// @Failure      401 {object} utils.AppError
// @Failure      422 {object} utils.AppError
// @Failure      500 {object} utils.AppError
// @Security     BearerAuth
// @Router       /tasks [post]
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

	taskJSON, _ := json.Marshal(task)
	h.hub.Broadcast(task.BoardID.Hex(), "task_created", string(taskJSON))

	return utils.WriteJSON(w, http.StatusCreated, task)
}

// UpdateTask godoc
// @Summary      Update a task
// @Description  Update one or more fields of an existing task. Only provided fields are updated.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path     string           true  "Task ID"
// @Param        task body     models.UpdateTask true  "Fields to update"
// @Success      200  {object} map[string]string
// @Failure      400  {object} utils.AppError
// @Failure      401  {object} utils.AppError
// @Failure      404  {object} utils.AppError
// @Failure      500  {object} utils.AppError
// @Security     BearerAuth
// @Router       /tasks/{id} [patch]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	var req models.UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.BadRequest("invalid request body")
	}

	if err := h.service.UpdateTask(r.Context(), id, &req); err != nil {
		return err
	}

	reqJSON, _ := json.Marshal(req)
	h.hub.Broadcast(id, "task_updated", string(reqJSON))

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "task updated successfully",
	})
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Permanently delete a task by ID
// @Tags         tasks
// @Produce      json
// @Param        id  path     string true "Task ID"
// @Success      200 {object} map[string]string
// @Failure      401 {object} utils.AppError
// @Failure      404 {object} utils.AppError
// @Failure      500 {object} utils.AppError
// @Security     BearerAuth
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		return err
	}

	h.hub.Broadcast(id, "task_deleted", fmt.Sprintf(`{"id": "%s"}`, id))

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "task deleted successfully",
	})
}

// GetTasks godoc
// @Summary      Get all tasks
// @Description  Get paginated tasks for a board with optional filters
// @Tags         tasks
// @Produce      json
// @Param        board_id       query    string false "Board ID (required)"
// @Param        column_id      query    string false "Filter by column ID"
// @Param        destination_id query    string false "Filter by assignee ID"
// @Param        priority       query    string false "Filter by priority (low, medium, high)"
// @Param        search         query    string false "Search in title and description"
// @Param        page           query    int    false "Page number (default 1)"
// @Param        limit          query    int    false "Items per page (default 10, max 100)"
// @Success      200 {object}   models.PaginatedTasks
// @Failure      400 {object}   utils.AppError
// @Failure      401 {object}   utils.AppError
// @Failure      500 {object}   utils.AppError
// @Security     BearerAuth
// @Router       /tasks [get]
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) error {
	boardID := r.URL.Query().Get("board_id")
	columnID := r.URL.Query().Get("column_id")
	destinationID := r.URL.Query().Get("destination_id")
	priority := r.URL.Query().Get("priority")
	search := r.URL.Query().Get("search")

	page := int64(1)
	if p := r.URL.Query().Get("page"); p != "" {
		parsed, err := strconv.ParseInt(p, 10, 64)
		if err != nil || parsed < 1 {
			return utils.BadRequest("page must be a positive number")
		}
		page = parsed
	}

	limit := int64(10)
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.ParseInt(l, 10, 64)
		if err != nil || parsed < 1 {
			return utils.BadRequest("limit must be a positive number")
		}
		limit = parsed
	}

	filter := &models.TaskFilter{
		BoardID:       boardID,
		ColumnID:      columnID,
		DestinationID: destinationID,
		Priority:      priority,
		Search:        search,
		Page:          page,
		Limit:         limit,
	}

	tasks, err := h.service.GetTasks(r.Context(), filter)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, tasks)
}

// GetTask godoc
// @Summary      Get a task
// @Description  Get a single task by ID with column details populated
// @Tags         tasks
// @Produce      json
// @Param        id  path     string true "Task ID"
// @Success      200 {object} models.TaskResponse
// @Failure      400 {object} utils.AppError
// @Failure      401 {object} utils.AppError
// @Failure      404 {object} utils.AppError
// @Failure      500 {object} utils.AppError
// @Security     BearerAuth
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, task)
}

// CompleteTask godoc
// @Summary      Update task column
// @Description  Move a task to a different column
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path     string        true "Task ID"
// @Param        body body     models.Column true "Column object"
// @Success      200  {object} map[string]string
// @Failure      400  {object} utils.AppError
// @Failure      401  {object} utils.AppError
// @Failure      404  {object} utils.AppError
// @Security     BearerAuth
// @Router       /tasks/{id}/complete [patch]
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

// UpdateColumn godoc
// @Summary      Update a column
// @Description  Update the name of a column. Valid names are todo, in progress, completed. All tasks in the column automatically reflect the new name via $lookup.
// @Tags         columns
// @Accept       json
// @Produce      json
// @Param        id     path     string             true "Column ID"
// @Param        column body     models.UpdateColumn true "Column update object"
// @Success      200    {object} map[string]string
// @Failure      400    {object} utils.AppError
// @Failure      401    {object} utils.AppError
// @Failure      404    {object} utils.AppError
// @Failure      500    {object} utils.AppError
// @Security     BearerAuth
// @Router       /columns/{id} [patch]
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

// CheckId godoc
// @Summary      Check ID
// @Description  Validate a numeric ID from the URL path
// @Tags         utils
// @Produce      plain
// @Param        id  path     int true "Numeric ID"
// @Success      200 {string} string
// @Failure      404 {object} utils.AppError
// @Security     BearerAuth
// @Router       /check/{id} [get]
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
