package repository

import (
	"context"
	"fmt"
	"task-tracker-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BoardRepository struct {
	boardCollection  *mongo.Collection
	columnCollection *mongo.Collection
	taskCollection   *mongo.Collection
}

func NewBoardReposity(boardCollection *mongo.Collection, columnCollection *mongo.Collection, taskCollection *mongo.Collection) *BoardRepository {
	return &BoardRepository{
		boardCollection:  boardCollection,
		columnCollection: columnCollection,
		taskCollection:   taskCollection,
	}
}

func (r *BoardRepository) CreateBoard(ctx context.Context, req *models.Board) error {

	req.ID = primitive.NewObjectID()

	_, err := r.boardCollection.InsertOne(ctx, req)

	if err != nil {
		return fmt.Errorf("Error inserting board")
	}

	return nil

}

// repository/board_repository.go
func (r *BoardRepository) GetBoard(ctx context.Context, id string) (*models.Board, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid board id")
	}

	var board models.Board
	err = r.boardCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&board)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("board not found")
		}
		return nil, fmt.Errorf("error fetching board")
	}

	return &board, nil
}

// repository/board_repository.go
func (r *BoardRepository) GetAllBoards(ctx context.Context) ([]models.Board, error) {
	var boards []models.Board

	cursor, err := r.boardCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching boards")
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &boards); err != nil {
		return nil, fmt.Errorf("error decoding boards")
	}

	return boards, nil
}

// repository/board_repository.go
func (r *BoardRepository) UpdateBoard(ctx context.Context, id string, req *models.UpdateBoard) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid board id")
	}

	fields := bson.M{}

	if req.Name != nil && *req.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if req.Name != nil {
		fields["name"] = *req.Name
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	result, err := r.boardCollection.UpdateOne(ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": fields},
	)
	if err != nil {
		return fmt.Errorf("error updating board")
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("board not found")
	}

	return nil
}

// repository/board_repository.go
func (r *BoardRepository) GetBoardsByDestination(ctx context.Context, destinationID string) ([]models.Board, error) {
	var boards []models.Board

	cursor, err := r.boardCollection.Find(ctx, bson.M{"destination_id": destinationID})
	if err != nil {
		return nil, fmt.Errorf("error fetching boards")
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &boards); err != nil {
		return nil, fmt.Errorf("error decoding boards")
	}

	return boards, nil
}

// repository/board_repository.go
func (r *BoardRepository) DeleteBoard(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid board id")
	}

	session, err := r.boardCollection.Database().Client().StartSession()
	if err != nil {
		return fmt.Errorf("error starting session")
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		// 1. delete board
		result, err := r.boardCollection.DeleteOne(sessCtx, bson.M{"_id": objectID})
		if err != nil {
			return nil, fmt.Errorf("error deleting board")
		}
		if result.DeletedCount == 0 {
			return nil, fmt.Errorf("board not found")
		}

		// 2. delete all columns in this board
		_, err = r.columnCollection.DeleteMany(sessCtx, bson.M{"board_id": objectID})
		if err != nil {
			return nil, fmt.Errorf("error deleting columns")
		}

		// 3. delete all tasks in this board
		_, err = r.taskCollection.DeleteMany(sessCtx, bson.M{"board_id": objectID})
		if err != nil {
			return nil, fmt.Errorf("error deleting tasks")
		}

		return nil, nil
	})

	return err
}
