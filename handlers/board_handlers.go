package handlers

import (
	"encoding/json"
	"net/http"
	"task-tracker-api/models"
	"task-tracker-api/services"
	"task-tracker-api/utils"
)

type BoardHandler struct {
	service *services.BoardServices
}

func NewBoardHandler(service *services.BoardServices) *BoardHandler {
	return &BoardHandler{service: service}
}

// CreateBoard godoc
// @Summary      Create board
// @Description  Create a new board under a destination
// @Tags         boards
// @Accept       json
// @Produce      json
// @Param        id   path      string        true  "Destination ID"
// @Param        board body     models.Board  true  "Board payload"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  utils.AppError
// @Failure      500  {object}  utils.AppError
// @Router       /destinations/{id}/boards [post]
func (h *BoardHandler) CreateBoard(w http.ResponseWriter, r *http.Request) error {
	var req models.Board

	req.DestinationID = r.PathValue("id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.BadRequest(err.Error())
	}

	err := h.service.CreateBoard(r.Context(), &req)
	if err != nil {
		return utils.InternalServerError(err.Error())
	}

	return utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Board created",
	})
}

// GetBoard godoc
// @Summary      Get board
// @Description  Retrieve a single board by ID
// @Tags         boards
// @Produce      json
// @Param        id   path      string  true  "Board ID"
// @Success      200  {object}  models.Board
// @Failure      404  {object}  utils.AppError
// @Failure      500  {object}  utils.AppError
// @Router       /boards/{id} [get]
// handlers/board_handler.go
func (h *BoardHandler) GetBoard(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	board, err := h.service.GetBoard(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, board)
}

// GetAllBoards godoc
// @Summary      Get all boards
// @Description  Retrieve all boards
// @Tags         boards
// @Produce      json
// @Success      200  {array}   models.Board
// @Failure      500  {object}  utils.AppError
// @Router       /boards [get]
// handlers/board_handler.go
func (h *BoardHandler) GetAllBoards(w http.ResponseWriter, r *http.Request) error {
	boards, err := h.service.GetAllBoards(r.Context())
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, boards)
}

// UpdateBoard godoc
// @Summary      Update board
// @Description  Update a board by ID
// @Tags         boards
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "Board ID"
// @Param        board body      models.UpdateBoard   true  "Updated board payload"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  utils.AppError
// @Failure      404   {object}  utils.AppError
// @Failure      500   {object}  utils.AppError
// @Router       /boards/{id} [put]
// handlers/board_handler.go
func (h *BoardHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	var req models.UpdateBoard
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.BadRequest("invalid request body")
	}

	if err := h.service.UpdateBoard(r.Context(), id, &req); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "board updated successfully",
	})
}

// GetBoardsByDestination godoc
// @Summary      Get boards by destination
// @Description  Retrieve boards filtered by destination ID
// @Tags         boards
// @Produce      json
// @Param        destination_id  query     string  true  "Destination ID"
// @Success      200             {array}   models.Board
// @Failure      400             {object}  utils.AppError
// @Failure      500             {object}  utils.AppError
// @Router       /boards/by-destination [get]
// handlers/board_handler.go
func (h *BoardHandler) GetBoardsByDestination(w http.ResponseWriter, r *http.Request) error {
	destinationID := r.URL.Query().Get("destination_id")

	boards, err := h.service.GetBoardsByDestination(r.Context(), destinationID)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, boards)
}

// DeleteBoard godoc
// @Summary      Delete board
// @Description  Delete a board by ID
// @Tags         boards
// @Produce      json
// @Param        id   path      string  true  "Board ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  utils.AppError
// @Failure      500  {object}  utils.AppError
// @Router       /boards/{id} [delete]
// handlers/board_handler.go
func (h *BoardHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	if err := h.service.DeleteBoard(r.Context(), id); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "board deleted successfully",
	})
}
