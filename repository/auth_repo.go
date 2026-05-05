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

func (r *AuthRepository) GetUser(ctx context.Context, userID string) (*models.UserResponse, error) {
	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	pipeline := bson.A{
		// match the user
		bson.M{"$match": bson.M{"_id": objectID}},

		// convert team_id strings to ObjectIDs for lookup
		bson.M{"$addFields": bson.M{
			"team_object_ids": bson.M{
				"$map": bson.M{
					"input": "$team_id",
					"as":    "tid",
					"in": bson.M{
						"$toObjectId": "$$tid",
					},
				},
			},
		}},

		// lookup teams
		bson.M{"$lookup": bson.M{
			"from":         "teams",
			"localField":   "team_object_ids",
			"foreignField": "_id",
			"as":           "teams",
		}},

		// remove internal fields
		bson.M{"$project": bson.M{
			"password":        0,
			"team_id":         0,
			"team_object_ids": 0,
		}},
	}

	cursor, err := r.userCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error fetching user")
	}
	defer cursor.Close(ctx)

	var results []models.UserResponse
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("error decoding user")
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &results[0], nil
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
