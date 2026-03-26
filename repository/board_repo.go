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
	boardCollection *mongo.Collection
}

func NewBoardReposity(boardCollection *mongo.Collection) *BoardRepository {
	return &BoardRepository{
		boardCollection: boardCollection,
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
