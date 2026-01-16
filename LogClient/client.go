package logClient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/EsanSamuel/sensory/db"
	"github.com/EsanSamuel/sensory/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Client struct {
	conn    net.Conn
	Project string
	Service string
	ApiKey  string
}

type LogEntry struct {
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	Project   string `json:"project"`
	Service   string `json:"service"`
	Message   string `json:"message"`
}

func New(apikey, addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Warning: Failed to connect to log server: %v", err)
		return nil, err
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var project models.Project

	err = db.ProjectCollection.FindOne(ctx, bson.M{"api_key": apikey}).Decode(&project)
	if err != nil {
		log.Println("Cannot find project", err)
		return nil, err
	}

	return &Client{conn: conn, Project: project.ProjectName, Service: project.Service, ApiKey: apikey}, nil
}

func (c *Client) Send(level string, msg string) error {
	entry := LogEntry{
		Level:     level,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Project:   c.Project,
		Service:   c.Service,
		Message:   msg,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Println(err)
	}

	data = append(data, '\n')
	_, err = c.conn.Write(data)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
