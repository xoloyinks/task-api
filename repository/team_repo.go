package repository

import (
	"context"
	"fmt"
	"task-tracker-api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TeamRepository struct {
	collection           *mongo.Collection
	teamMemberCollection *mongo.Collection
	userCollection       *mongo.Collection
}

func NewTeamRepository(collection *mongo.Collection, userCollection *mongo.Collection, teamMemberCollection *mongo.Collection) *TeamRepository {
	return &TeamRepository{
		collection:           collection,
		userCollection:       userCollection,
		teamMemberCollection: teamMemberCollection,
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

func (r *TeamRepository) AddMember(ctx context.Context, teamID string, email string, role string) error {

	session, err := r.userCollection.Database().Client().StartSession()
	if err != nil {
		return fmt.Errorf("error starting session")
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		// operation 1
		filter := bson.M{"email": email}
		update := bson.M{
			"$push": bson.M{
				"team_id": teamID,
			},
		}

		result, err := r.userCollection.UpdateOne(sessCtx, filter, update)
		if err != nil {
			return nil, fmt.Errorf("error adding member to team")
		}

		if result.MatchedCount == 0 {
			return nil, fmt.Errorf("user not found")
		}

		// operation 2
		member := models.TeamMember{
			ID:        primitive.NewObjectID(),
			TeamID:    teamID,
			Email:     email,
			Role:      role,
			CreatedAt: time.Now(),
		}

		_, err = r.teamMemberCollection.InsertOne(sessCtx, member)
		if err != nil {
			return nil, fmt.Errorf("error creating team member")
		}

		return nil, nil
	})

	return err
}

func (r *TeamRepository) GetAllTeams(ctx context.Context) ([]models.Team, error) {
	var teams []models.Team
	cursor, err := r.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, fmt.Errorf("Error fetching docs")
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &teams); err != nil {
		return nil, err
	}

	return teams, nil
}

func (r *TeamRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
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
