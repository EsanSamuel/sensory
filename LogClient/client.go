package logClient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"runtime"
	"time"

	"github.com/EsanSamuel/sensory/db"
	"github.com/EsanSamuel/sensory/helpers"
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
	noOp      bool // flag to mark dummy client
}

type Runtime struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Fn   string `json:"fn"`
}

type LogEntry struct {
	Level     string  `json:"level"`
	Timestamp string  `json:"timestamp"`
	Project   string  `json:"project"`
	Service   string  `json:"service"`
	Message   string  `json:"message"`
	Runtime   Runtime `json:"runtime"`
}

func New(apikey, addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var project models.Project
	err = db.ProjectCollection.FindOne(ctx, bson.M{"api_key": apikey}).Decode(&project)
	if err != nil {
		conn.Close()
		return nil, err
	}

	client := &Client{
		conn:      conn,
		Project:   project.ProjectName,
		Service:   project.Service,
		ApiKey:    apikey,
		ProjectId: project.ProjectID,
		UserId:    project.UserID,
		noOp:      false,
	}

	err = helpers.SaveProjectMeta(client.Project, client.ProjectId, client.UserId)
	if err != nil {
		log.Println(err)
	}

	return client, nil
}

func NewNoOp() *Client {
	return &Client{noOp: true}
}

func getLocation() (file string, line int, function string) {
	pc, file, line, ok := runtime.Caller(3)

	if !ok {
		return "unknown", 0, "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file, line, "unknown"
	}
	return file, line, fn.Name()
}

func (c *Client) Send(level, msg string) error {
	if c == nil || c.noOp {
		fmt.Println("[noop logger]", level, msg)
		return nil
	}

	file, line, fn := getLocation()

	entry := LogEntry{
		Level:     level,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Project:   c.Project,
		Service:   c.Service,
		Message:   msg,
		Runtime: Runtime{
			File: file,
			Line: line,
			Fn:   fn,
		},
	}

	PushLogToDB(entry, c)

	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return err
	}

	data = append(data, '\n')
	if c.conn != nil {
		_, err = c.conn.Write(data)
		if err != nil {
			fmt.Println("TCP write failed:", err)
		}
	}
	return err
}

func PushLogToDB(entry LogEntry, c *Client) {
	if c == nil || c.noOp {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	logEntry := models.Log{
		LogID:     bson.NewObjectID().Hex(),
		Level:     entry.Level,
		Message:   entry.Message,
		Service:   entry.Service,
		TimeStamp: entry.Timestamp,
		ProjectID: c.ProjectId,
		UserID:    c.UserId,
		Runtime: models.Runtime{
			File: entry.Runtime.File,
			Line: entry.Runtime.Line,
			Fn:   entry.Runtime.Fn,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := db.LogCollection.InsertOne(ctx, logEntry)
	if err != nil {
		fmt.Println("DB insert error:", err)
		return
	}

	fmt.Println("saved log:", result)
}
