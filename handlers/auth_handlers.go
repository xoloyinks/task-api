package handlers

import (
	"encoding/json"
	"net/http"
	"task-tracker-api/models"
	"task-tracker-api/services"
	"task-tracker-api/utils"
)

type AuthHandler struct {
	service *services.AuthServices
}

func NewAuthHandler(service *services.AuthServices) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) CreateAccount(w http.ResponseWriter, r *http.Request) error {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(err.Error())
	}

	if err := h.service.CreateAccount(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Account created!"})

	return nil
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req models.Login

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return utils.BadRequest(err.Error())
	}
	user, token, err := h.service.Login(r.Context(), &req)

	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, models.ActiveUser{
		Token: token,
		User: models.ActiveUserResponse{
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreateAt,
			ID:        user.ID,
		},
	})
}
