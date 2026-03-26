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

type TaskRepository struct {
	taskCollection   *mongo.Collection
	columnCollection *mongo.Collection
	boardCollection  *mongo.Collection
}

func NewTaskRepository(taskCollection *mongo.Collection, columnCollection *mongo.Collection, boardCollection *mongo.Collection) *TaskRepository {
	return &TaskRepository{
		taskCollection:   taskCollection,
		columnCollection: columnCollection,
		boardCollection:  boardCollection,
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	session, err := r.taskCollection.Database().Client().StartSession()
	if err != nil {
		return fmt.Errorf("error starting session")
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (interface{}, error) {
		// 1. check board exists
		boardObjectID, err := primitive.ObjectIDFromHex(task.BoardID.Hex())
		if err != nil {
			return nil, fmt.Errorf("invalid board id")
		}

		var board models.Board
		err = r.boardCollection.FindOne(sessCtx, bson.M{"_id": boardObjectID}).Decode(&board)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("board not found")
			}
			return nil, fmt.Errorf("error finding board")
		}

		// 2. create default column
		column := models.Column{
			ID:      primitive.NewObjectID(),
			BoardID: task.BoardID,
			Name:    "todo", // default column name
		}

		_, err = r.columnCollection.InsertOne(sessCtx, column)
		if err != nil {
			return nil, fmt.Errorf("error creating column")
		}

		// 3. create task with column ID
		task.ID = primitive.NewObjectID()
		task.ColumnID = column.ID
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()

		_, err = r.taskCollection.InsertOne(sessCtx, task)
		if err != nil {
			return nil, fmt.Errorf("error creating task")
		}

		return nil, nil
	})

	return err
}

func (r *TaskRepository) GetAllTasks(ctx context.Context, userId string) (error, []models.Task) {
	var tasks []models.Task

	cursor, err := r.taskCollection.Find(ctx, bson.M{
		"user_id": userId,
	})
	if err != nil {
		return err, nil
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &tasks); err != nil {
		return err, nil
	}

	return nil, tasks
}

func (r *TaskRepository) GetTask(ctx context.Context, id string) (*models.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid task id")
	}

	filter := bson.M{"_id": objectID}

	var task models.Task
	result := r.taskCollection.FindOne(ctx, filter)
	if err := result.Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("error fetching task")
	}

	return &task, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, id string, req *models.UpdateTask) error {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return fmt.Errorf("Invalid ID")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"title":       req.Title,
			"description": req.Description,
			"updated_at":  time.Now(),
		},
	}

	result, err := r.taskCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Error updating task")
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("task not fount")
	}

	return nil

}

func (r *TaskRepository) Column(ctx context.Context, id string, req *models.Column) error {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return fmt.Errorf("Invalid task id")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"column":     req.Name,
			"updated_at": time.Now(),
		},
	}

	// mongodb update here
	result, err := r.taskCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Error updating task")
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("task not found")
	}

	return nil

}

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return fmt.Errorf("Invalid task id")
	}

	// delete task with id
	filter := bson.M{"_id": objectID}

	result, err := r.taskCollection.DeleteOne(ctx, filter)

	if err != nil {
		return fmt.Errorf("Error deleting task")
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("task not found")
	}

	return nil

}
