package logClient

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/EsanSamuel/sensory/db"
	"github.com/EsanSamuel/sensory/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Client struct {
	conn      net.Conn
	Project   string
	Service   string
	ApiKey    string
	ProjectId string
	UserId    string
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
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var project models.Project
	err = db.ProjectCollection.FindOne(
		ctx,
		bson.M{"api_key": apikey},
	).Decode(&project)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		conn:      conn,
		Project:   project.ProjectName,
		Service:   project.Service,
		ApiKey:    apikey,
		ProjectId: project.ProjectID,
		UserId:    project.UserID,
	}, nil
}

func (c *Client) Send(level string, msg string) error {
	entry := LogEntry{
		Level:     level,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Project:   c.Project,
		Service:   c.Service,
		Message:   msg,
	}

	PushLogToDB(entry, c)

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

func PushLogToDB(entry LogEntry, c *Client) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var log models.Log

	log.LogID = bson.NewObjectID().Hex()
	log.Level = entry.Level
	log.Message = entry.Message
	log.Service = entry.Service
	log.TimeStamp = entry.Timestamp
	log.ProjectID = c.ProjectId
	log.UserID = c.UserId
	log.CreatedAt = time.Now()
	log.UpdatedAt = time.Now()

	result, err := db.LogCollection.InsertOne(ctx, log)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("saved log: ", result)
}
