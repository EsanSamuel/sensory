package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Project struct {
	ID          bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	ProjectID   string        `json:"project_id" bson:"project_id"`
	ProjectName string        `json:"project_name" bson:"project_name" validate:"required,min=2,max=100"`
	Description string        `json:"description" bson:"description"`
	UserID      string        `json:"user_id" bson:"user_id"`
	ApiKey      string        `json:"api_key" bson:"api_key"`
	Service     string        `json:"service" bson:"service"`
	LogCounts   int           `json:"log_counts" bson:"log_counts"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

type Runtime struct {
	File string `json:"file" bson:"file"`
	Line int    `json:"line" bson:"line"`
	Fn   string `json:"fn" bson:"fn"`
}

type Log struct {
	ID        bson.ObjectID `json:"_id" bson:"_id,omitempty"`
	LogID     string        `json:"log_id" bson:"log_id"`
	ProjectID string        `json:"project_id" bson:"project_id"`
	UserID    string        `json:"user_id" bson:"user_id"`
	Level     string        `json:"level" bson:"level"`
	TimeStamp string        `json:"timestamp" bson:"timestamp"`
	Message   string        `json:"message" bson:"message"`
	Service   string        `json:"service" bson:"service"`
	Runtime   Runtime       `json:"runtime" bson:"runtime"`
	Project   Project       `json:"project" bson:"project"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
