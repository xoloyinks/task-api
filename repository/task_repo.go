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
	collection     *mongo.Collection
	userCollection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection, userCollection *mongo.Collection) *TaskRepository {
	return &TaskRepository{
		collection:     collection,
		userCollection: userCollection,
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	task.ID = primitive.NewObjectID()
	task.CreartAt = time.Now()

	_, err := r.collection.InsertOne(ctx, task)

	return err
}

func (r *TaskRepository) GetAllTasks(ctx context.Context, userId string) (error, []models.Task) {
	var tasks []models.Task

	cursor, err := r.collection.Find(ctx, bson.M{
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
	result := r.collection.FindOne(ctx, filter)
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

	result, err := r.collection.UpdateOne(ctx, filter, update)
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
	result, err := r.collection.UpdateOne(ctx, filter, update)
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

	result, err := r.collection.DeleteOne(ctx, filter)

	if err != nil {
		return fmt.Errorf("Error deleting task")
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("task not found")
	}

	return nil

}
