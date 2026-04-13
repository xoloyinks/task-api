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

func taskPipeline(match bson.M) bson.A {
	return bson.A{
		bson.M{"$match": match},
		bson.M{"$lookup": bson.M{
			"from":         "columns",
			"localField":   "column_id",
			"foreignField": "_id",
			"as":           "column",
		}},
		bson.M{"$unwind": bson.M{
			"path":                       "$column",
			"preserveNullAndEmptyArrays": true, // dont fail if column is missing
		}},
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
		var board models.Board
		err = r.boardCollection.FindOne(sessCtx, bson.M{"_id": task.BoardID}).Decode(&board)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("board not found")
			}
			return nil, fmt.Errorf("error finding board")
		}

		// 2. check if todo column already exists for this board
		var column models.Column
		err = r.columnCollection.FindOne(sessCtx, bson.M{
			"board_id": task.BoardID,
			"name":     "todo",
		}).Decode(&column)

		if err == mongo.ErrNoDocuments {
			// column doesn't exist — create it
			column = models.Column{
				ID:      bson.NewObjectID(),
				BoardID: task.BoardID,
				Name:    "todo",
			}
			_, err = r.columnCollection.InsertOne(sessCtx, column)
			if err != nil {
				return nil, fmt.Errorf("error creating column")
			}
		} else if err != nil {
			return nil, fmt.Errorf("error finding column")
		}

		// 3. create task with existing or new column ID
		task.ID = bson.NewObjectID()
		task.ColumnID = column.ID // reuse existing column
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

func (r *TaskRepository) GetTasks(ctx context.Context, filter *models.TaskFilter) (*models.PaginatedTasks, error) {
	// 1. build match filter
	match := bson.M{}

	// board_id — required
	boardObjectID, err := primitive.ObjectIDFromHex(filter.BoardID)
	if err != nil {
		return nil, fmt.Errorf("invalid board id")
	}
	match["board_id"] = boardObjectID

	// column_id — optional
	if filter.ColumnID != "" {
		columnObjectID, err := primitive.ObjectIDFromHex(filter.ColumnID)
		if err != nil {
			return nil, fmt.Errorf("invalid column id")
		}
		match["column_id"] = columnObjectID
	}

	// destination_id — optional
	if filter.DestinationID != "" {
		match["destination_id"] = filter.DestinationID
	}

	// priority — optional
	if filter.Priority != "" {
		match["priority"] = filter.Priority
	}

	// search — optional, searches title and description
	if filter.Search != "" {
		match["$or"] = bson.A{
			bson.M{"title": bson.M{"$regex": filter.Search, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// 2. calculate skip
	skip := (filter.Page - 1) * filter.Limit

	// 3. count total matching documents
	countPipeline := bson.A{
		bson.M{"$match": match},
		bson.M{"$count": "total"},
	}

	countCursor, err := r.taskCollection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, fmt.Errorf("error counting tasks")
	}
	defer countCursor.Close(ctx)

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	countCursor.All(ctx, &countResult)

	total := int64(0)
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// 4. build full pipeline with pagination
	pipeline := bson.A{
		bson.M{"$match": match},
		bson.M{"$lookup": bson.M{
			"from":         "columns",
			"localField":   "column_id",
			"foreignField": "_id",
			"as":           "column",
		}},
		bson.M{"$unwind": bson.M{
			"path":                       "$column",
			"preserveNullAndEmptyArrays": true,
		}},
		bson.M{"$sort": bson.M{"created_at": -1}}, // newest first
		bson.M{"$skip": skip},
		bson.M{"$limit": filter.Limit},
	}

	cursor, err := r.taskCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks")
	}
	defer cursor.Close(ctx)

	var tasks []models.TaskResponse
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("error decoding tasks")
	}

	// 5. calculate pagination metadata
	totalPages := (total + filter.Limit - 1) / filter.Limit

	return &models.PaginatedTasks{
		Data: tasks,
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    filter.Page < totalPages,
			HasPrev:    filter.Page > 1,
		},
	}, nil
}

func (r *TaskRepository) GetTask(ctx context.Context, id string) (*models.TaskResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid task id")
	}

	cursor, err := r.taskCollection.Aggregate(ctx, taskPipeline(bson.M{"_id": objectID}))
	if err != nil {
		return nil, fmt.Errorf("error fetching task")
	}
	defer cursor.Close(ctx)

	var tasks []models.TaskResponse
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("error decoding task")
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("task not found")
	}

	return &tasks[0], nil
}

// repository/task_repository.go
// repository/task_repository.go — convert string to ObjectID before storing
func (r *TaskRepository) UpdateTask(ctx context.Context, id string, req *models.UpdateTask) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid task id")
	}

	fields := bson.M{"updated_at": time.Now()}

	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.Description != nil {
		fields["description"] = *req.Description
	}
	if req.Priority != nil {
		fields["priority"] = *req.Priority
	}
	if req.Category != nil {
		fields["category"] = *req.Category
	}
	if req.DueDate != nil {
		fields["due_date"] = *req.DueDate
	}
	// ✅ convert string column_id to ObjectID
	if req.ColumnID != nil && *req.ColumnID != "" {
		colObjectID, err := primitive.ObjectIDFromHex(*req.ColumnID)
		if err != nil {
			return fmt.Errorf("invalid column_id")
		}
		fields["column_id"] = colObjectID
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": fields}

	result, err := r.taskCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating task")
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *TaskRepository) CreateColumn(ctx context.Context, col *models.Column) error {
	// check if column already exists for this board
	var existing models.Column
	err := r.columnCollection.FindOne(ctx, bson.M{
		"board_id": col.BoardID,
		"name":     col.Name,
	}).Decode(&existing)

	if err == nil {
		// already exists — return existing ID back to caller
		col.ID = existing.ID
		return nil
	}

	if err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking column")
	}

	// create it
	col.ID = bson.NewObjectID()
	_, err = r.columnCollection.InsertOne(ctx, col)
	if err != nil {
		return fmt.Errorf("error creating column")
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

func (r *TaskRepository) UpdateColumn(ctx context.Context, id string, req *models.UpdateColumn) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid column id")
	}

	result, err := r.columnCollection.UpdateOne(ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"name": *req.Name}},
	)
	if err != nil {
		return fmt.Errorf("error updating column")
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("column not found")
	}

	return nil
}
