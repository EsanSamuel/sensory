package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID    string        `json:"user_id" bson:"user_id"`
	Username  string        `json:"username" bson:"username"`
	Email     string        `json:"email" bson:"email" validate:"required,email"`
	Avatar    string        `json:"avatar" bson:"avatar"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
