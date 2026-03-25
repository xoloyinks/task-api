package repository

import (
	"context"
	"task-tracker-api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// create team
// add existing user to team

type TeamRepository struct {
	collection     *mongo.Collection
	userCollection *mongo.Collection
}

func NewTeamRepository(collection *mongo.Collection, userCollection *mongo.Collection) *TeamRepository {
	return &TeamRepository{
		collection:     collection,
		userCollection: userCollection,
	}
}

// create team
func (r *TeamRepository) CreateTeam(ctx context.Context, req *models.Team) error {

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, req)

	return err
}

func (r *TeamRepository) GetTeamByName(ctx context.Context, name string) (*models.Team, error) {
	var team models.Team
	err := r.userCollection.FindOne(ctx, bson.M{"name": name}).Decode(&team)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}
