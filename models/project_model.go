package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Project struct {
	ID          bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	ProjectName string        `json:"project_name" bson:"project_name" validate:"required,min=2,max=100"`
	UserID      string        `json:"user_id" bson:"user_id"`
	ApiKey      string        `json:"api_key" bson:"api_key"`
	Service     string        `json:"service" bson:"service"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

type Log struct {
	ID        bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	ProjectID string        `json:"project" bson:"project"`
	UserID    string        `json:"user_id" bson:"user_id"`
	Level     string        `json:"level" bson:"level"`
	TimeStamp string        `json:"timestamp" bson:"timestamp"`
	Message   string        `json:"message" bson:"message"`
	Service   string        `json:"service" bson:"service"`
}
