package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	TeamID   []string           `bson:"team_id" json:"team_id"`
	Password string             `bson:"password" json:"password"`
	CreateAt time.Time          `bson:"created_at" json:"created_at"`
}

type Team struct {
	Name        string             `bson:"name" json:"name"`
	ID          primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Description string             `bson:"description" json:"description"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type TeamMember struct {
	ID     primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	TeamID uint               `bson:"team_id" json:"team_id"`
	UserID uint               `bson:"user_id" json:"user_id"`
	Role   string             `bson:"role" json:"role"`
}

type Board struct {
	ID     primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	TeamID uint               `bson:"team_id" json:"team_id"`
	Name   string             `bson:"name" json:"name"`
}

type Column struct {
	ID      primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	BoardID uint               `bson:"board_id" json:"board_id"`
	Name    string             `bson:"name" json:"name"` // todo, in progress, completed
}

type Task struct {
	ID          primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	TeamID      string             `bson:"team_id" json:"team_id"`
	BoardID     uint               `bson:"board_id" json:"board_id"`
	ColumnID    uint               `bson:"column_id" json:"column_id"`
	Priority    string             `bson:"priority" json:"priority"`
	Category    string             `bson:"category" json:"category"`
	DueDate     time.Time          `bson:"due_date" json:"due_date"`
	CreartAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type UpdateTask struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `bson:"priority" json:"priority"`
	Category    string    `bson:"category" json:"category"`
	DueDate     time.Time `bson:"due_date" json:"due_date"`
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
