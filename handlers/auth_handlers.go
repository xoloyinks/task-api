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

// CreateAccount godoc
// @Summary      Create account
// @Description  Register a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User registration payload"
// @Success      201 {object} map[string]string "Account created successfully"
// @Failure      400 {object} utils.AppError
// @Failure      500 {object} utils.AppError
// @Router       /register [post]
func (h *AuthHandler) CreateAccount(w http.ResponseWriter, r *http.Request) error {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.BadRequest(err.Error())
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

// Login godoc
// @Summary      Login
// @Description  Login with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.Login true "Login credentials"
// @Success      200 {object} models.ActiveUser
// @Failure      400 {object} utils.AppError
// @Router       /login [post]
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
