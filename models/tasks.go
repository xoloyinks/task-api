package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	TeamID   []string           `bson:"team_id" json:"team_id"`
	Password string             `bson:"password" json:"password"`
	CreateAt time.Time          `bson:"created_at" json:"created_at"`
}

type Team struct {
	Name        string             `bson:"name" json:"name"`
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Description string             `bson:"description" json:"description"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type TeamResponse struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	Members     []TeamMember       `bson:"members" json:"members"`
}

// models/team.go
type UpdateTeam struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
type TeamMember struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TeamID    string             `bson:"team_id" json:"team_id"`
	Email     string             `bson:"email" json:"email"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// destination can be to team or personal
type Board struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DestinationID string             `bson:"destination_id" json:"destination_id"`
	Name          string             `bson:"name" json:"name"`
}

type Column struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BoardID primitive.ObjectID `bson:"board_id" json:"board_id"`
	Name    string             `bson:"name" json:"name"` // todo, in progress, completed
}

// tasks with team id are only displayed in team boards, while task with no team id are displayed only to the user boards
type Task struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	DestinationID string             `bson:"destination_id" json:"destination_id"`
	BoardID       primitive.ObjectID `bson:"board_id" json:"board_id"`
	ColumnID      primitive.ObjectID `bson:"column_id" json:"column_id"`
	Priority      string             `bson:"priority" json:"priority"`
	Category      string             `bson:"category" json:"category"`
	DueDate       time.Time          `bson:"due_date" json:"due_date"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type TaskResponse struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	DestinationID string             `bson:"destination_id" json:"destination_id"`
	BoardID       primitive.ObjectID `bson:"board_id" json:"board_id"`
	Priority      string             `bson:"priority" json:"priority"`
	Category      string             `bson:"category" json:"category"`
	DueDate       time.Time          `bson:"due_date" json:"due_date"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	Column        Column             `bson:"column" json:"column"` // populated via $lookup
}

// models/task.go
type UpdateTask struct {
	BoardID     *primitive.ObjectID `json:"board_id"` // client sends this
	Title       *string             `json:"title"`
	Description *string             `json:"description"`
	ColumnID    *primitive.ObjectID `json:"column_id"`
	Priority    *string             `json:"priority"`
	Category    *string             `json:"category"`
	DueDate     *time.Time          `json:"due_date"`
}

type Login struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

// models/user.go
type ActiveUserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"created_at"`
}

type ActiveUser struct {
	Token string             `json:"token"`
	User  ActiveUserResponse `json:"user"`
}

type UpdateColumn struct {
	Name *string `json:"name" validate:"required"`
}

// models/board.go
type UpdateBoard struct {
	Name *string `json:"name"`
}

type TaskFilter struct {
	BoardID       string `json:"board_id"`
	ColumnID      string `json:"column_id"`
	DestinationID string `json:"destination_id"`
	Priority      string `json:"priority"`
	Search        string `json:"search"`
	Page          int64  `json:"page"`
	Limit         int64  `json:"limit"`
}

type PaginatedTasks struct {
	Data       []TaskResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

type Pagination struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}
