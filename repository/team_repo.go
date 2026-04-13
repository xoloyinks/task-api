package repository

import (
	"context"
	"fmt"
	"log"
	"task-tracker-api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TeamRepository struct {
	collection           *mongo.Collection
	userCollection       *mongo.Collection
	teamMemberCollection *mongo.Collection
	taskCollection       *mongo.Collection
	boardCollection      *mongo.Collection
	columnCollection     *mongo.Collection
}

func NewTeamRepository(collection *mongo.Collection, userCollection *mongo.Collection, teamMemberCollection *mongo.Collection, boardCollection *mongo.Collection, taskCollection *mongo.Collection, columnCollection *mongo.Collection) *TeamRepository {
	return &TeamRepository{
		collection:           collection,
		userCollection:       userCollection,
		teamMemberCollection: teamMemberCollection,
		taskCollection:       taskCollection,
		boardCollection:      boardCollection,
		columnCollection:     columnCollection,
	}
}

// create team
func (r *TeamRepository) CreateTeam(ctx context.Context, req *models.Team) (string, error) {

	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return "", fmt.Errorf("error starting session")
	}
	defer session.EndSession(ctx)

	req.ID = bson.NewObjectID()
	teamIDStr := req.ID.Hex()

	log.Printf("Creating team with ID: %s", req.ID)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {

		// operation 1 - create team

		req.CreatedAt = time.Now()

		_, err := r.collection.InsertOne(sessCtx, req)
		if err != nil {
			return nil, fmt.Errorf("error inserting team: %w", err)
		}

		// operation 2 - add teamID to user's team_id slice
		result, err := r.userCollection.UpdateOne(
			sessCtx,
			bson.M{"email": req.CreatedBy},
			bson.M{
				"$addToSet": bson.M{
					"team_id": teamIDStr,
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error updating user: %w", err)
		}

		if result.MatchedCount == 0 {
			return nil, fmt.Errorf("creator user not found")
		}

		// operation 3 - create team member document for creator with admin role
		member := models.TeamMember{
			ID:        bson.NewObjectID(),
			TeamID:    teamIDStr,
			Email:     req.CreatedBy,
			Role:      "admin",
			CreatedAt: time.Now(),
		}
		_, err = r.teamMemberCollection.InsertOne(sessCtx, member)
		if err != nil {
			return nil, fmt.Errorf("error creating team member: %w", err)
		}

		return nil, nil
	})

	return teamIDStr, err
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
			return nil, fmt.Errorf("error adding member to team: %w", err)
		}

		if result.MatchedCount == 0 {
			return nil, fmt.Errorf("user not found")
		}

		// operation 2
		member := models.TeamMember{
			ID:        bson.NewObjectID(),
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

// repository/team_repository.go
func (r *TeamRepository) GetTeams(ctx context.Context, email string) ([]models.Team, error) {
	var teams []models.Team

	cursor, err := r.collection.Find(ctx, bson.M{"created_by": email})
	if err != nil {
		return nil, fmt.Errorf("error fetching teams")
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &teams); err != nil {
		return nil, fmt.Errorf("error decoding teams")
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

// repository/team_repository.go
// richer response — team with its members
func (r *TeamRepository) GetTeam(ctx context.Context, id string) (*models.TeamResponse, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid team id")
	}

	log.Printf("Fetching team with ID: %s", objectID)

	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": objectID}},
		bson.M{
			"$lookup": bson.M{
				"from": "members",
				"let":  bson.M{"teamId": "$_id"},
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": bson.A{
									bson.M{"$toObjectId": "$team_id"}, // convert string → ObjectId
									"$$teamId",
								},
							},
						},
					},
				},
				"as": "members",
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error fetching team")
	}
	defer cursor.Close(ctx)

	var teams []models.TeamResponse
	if err := cursor.All(ctx, &teams); err != nil {
		return nil, fmt.Errorf("error decoding team")
	}

	if len(teams) == 0 {
		return nil, fmt.Errorf("team not found")
	}

	return &teams[0], nil
}

// repository/team_repository.go
func (r *TeamRepository) UpdateTeam(ctx context.Context, id string, req *models.UpdateTeam) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid team id")
	}

	fields := bson.M{}

	if req.Name != nil && *req.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if req.Name != nil {
		fields["name"] = *req.Name
	}

	if req.Description != nil {
		fields["description"] = *req.Description
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	result, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": fields},
	)
	if err != nil {
		return fmt.Errorf("error updating team")
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("team not found")
	}

	return nil
}

// repository/team_repository.go
func (r *TeamRepository) DeleteTeam(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid team id")
	}

	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return fmt.Errorf("error starting session")
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		// 1. check team exists
		result, err := r.collection.DeleteOne(sessCtx, bson.M{"_id": objectID})
		if err != nil {
			return nil, fmt.Errorf("error deleting team")
		}
		if result.DeletedCount == 0 {
			return nil, fmt.Errorf("team not found")
		}

		// 2. remove teamID from all affected users team_id slice
		_, err = r.userCollection.UpdateMany(sessCtx,
			bson.M{"team_id": id},
			bson.M{"$pull": bson.M{"team_id": id}},
		)
		if err != nil {
			return nil, fmt.Errorf("error updating users")
		}

		// 3. delete all members in that team
		_, err = r.teamMemberCollection.DeleteMany(sessCtx, bson.M{"team_id": id})
		if err != nil {
			return nil, fmt.Errorf("error deleting members")
		}

		// 4. get all boards in that team
		cursor, err := r.boardCollection.Find(sessCtx, bson.M{"destination_id": id})
		if err != nil {
			return nil, fmt.Errorf("error fetching boards")
		}
		defer cursor.Close(sessCtx)

		var boards []models.Board
		if err := cursor.All(sessCtx, &boards); err != nil {
			return nil, fmt.Errorf("error decoding boards")
		}

		// 5. delete all tasks and boards
		for _, board := range boards {
			// delete all tasks in each board
			_, err = r.taskCollection.DeleteMany(sessCtx, bson.M{"board_id": board.ID})
			if err != nil {
				return nil, fmt.Errorf("error deleting tasks")
			}

			// delete all columns in each board
			_, err = r.columnCollection.DeleteMany(sessCtx, bson.M{"board_id": board.ID})
			if err != nil {
				return nil, fmt.Errorf("error deleting columns")
			}
		}

		// 6. delete all boards in that team
		_, err = r.boardCollection.DeleteMany(sessCtx, bson.M{"destination_id": id})
		if err != nil {
			return nil, fmt.Errorf("error deleting boards")
		}

		return nil, nil
	})

	return err
}

// repository/team_repository.go
func (r *TeamRepository) RemoveMember(ctx context.Context, teamID string, memberID string) error {
	memberObjectID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return fmt.Errorf("invalid member id")
	}

	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return fmt.Errorf("error starting session")
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		// 1. find member before deleting — need email to update user
		var member models.TeamMember
		err := r.teamMemberCollection.FindOne(sessCtx, bson.M{
			"_id":     memberObjectID,
			"team_id": teamID,
		}).Decode(&member)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("member not found")
			}
			return nil, fmt.Errorf("error finding member")
		}

		// 2. delete member document
		_, err = r.teamMemberCollection.DeleteOne(sessCtx, bson.M{"_id": memberObjectID})
		if err != nil {
			return nil, fmt.Errorf("error removing member")
		}

		// 3. remove teamID from user's team_id slice
		_, err = r.userCollection.UpdateOne(sessCtx,
			bson.M{"email": member.Email},
			bson.M{"$pull": bson.M{"team_id": teamID}},
		)
		if err != nil {
			return nil, fmt.Errorf("error updating user")
		}

		return nil, nil
	})

	return err
}
