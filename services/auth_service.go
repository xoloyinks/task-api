package services

import (
	"context"
	"task-tracker-api/models"
	"task-tracker-api/repository"
	"task-tracker-api/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthServices struct {
	repo *repository.AuthRepository
}

func NewAuthServices(repo *repository.AuthRepository) *AuthServices {
	return &AuthServices{repo: repo}
}

func (s *AuthServices) CreateAccount(ctx context.Context, req *models.User) error {
	if err := validate.Struct(req); err != nil {
		return err
	}

	existingUser, _ := s.repo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return utils.BadRequest("email already in use")
	}

	if req.TeamID == nil {
		req.TeamID = []string{}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.InternalServerError("error creating account")
	}
	req.Password = string(hashedPassword)

	return s.repo.CreateAccount(ctx, req)
}

func (s *AuthServices) Login(ctx context.Context, req *models.Login) (*models.User, string, error) {
	// check if req is valid
	user, err := s.repo.Login(ctx, req)
	if err != nil {
		return nil, "", utils.BadRequest(err.Error())
	}

	// compare password with hashedpassword
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", utils.BadRequest("invalid email or password")
	}
	// generate jwt token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		return nil, "", utils.InternalServerError("error generating token")
	}

	return user, token, nil

}
