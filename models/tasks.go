package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	CreateAt time.Time          `bson:"created_at" json:"created_at"`
}

type Task struct {
	ID          primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	UserID      string             `bson:"user_id" json:"user_id"`
	Completed   bool               `bson:"completed" json:"completed"`
	CreartAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type UpdateTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CompleteTask struct {
	Completed bool `bson:"completed" json:"completed"`
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
