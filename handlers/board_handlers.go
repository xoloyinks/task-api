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

// handlers/board_handler.go
func (h *BoardHandler) GetBoard(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	board, err := h.service.GetBoard(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, board)
}

// handlers/board_handler.go
func (h *BoardHandler) GetAllBoards(w http.ResponseWriter, r *http.Request) error {
	boards, err := h.service.GetAllBoards(r.Context())
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, boards)
}
