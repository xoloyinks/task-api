package repository

import (
	"context"
	"fmt"
	"task-tracker-api/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepository struct {
	collection     *mongo.Collection
	userCollection *mongo.Collection
}

func NewAuthRepository(collection *mongo.Collection, userCollection *mongo.Collection) *AuthRepository {
	return &AuthRepository{
		collection:     collection,
		userCollection: userCollection,
	}
}

func (r *AuthRepository) CreateAccount(ctx context.Context, req *models.User) error {
	req.ID = bson.NewObjectID()
	req.CreateAt = time.Now()

	_, err := r.userCollection.InsertOne(ctx, req)

	if err != nil {
		return fmt.Errorf("Error creating user")
	}

	return nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// repository/auth_repository.go
func (r *AuthRepository) Login(ctx context.Context, req *models.Login) (*models.User, error) {
	var user models.User

	filter := bson.M{"email": req.Email}

	err := r.userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("error fetching user")
	}

	return &user, nil
}
