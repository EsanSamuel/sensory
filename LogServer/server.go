package logserver

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/EsanSamuel/sensory/db"
	"github.com/EsanSamuel/sensory/jobs/workers"
	"github.com/EsanSamuel/sensory/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Log struct {
	Level     string `json:"level"`
	TimeStamp string `json:"timestamp"`
	Message   string `json:"message"`
	Service   string `json:"service"`
	Project   string `json:"project"`
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
}

func Initialize_Log() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	l, err := net.Listen("tcp", ":10000")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer l.Close()

	fmt.Println("TCP Log Server listening on :9000")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connection from:", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()

		var entry models.LogEntry
		err := json.Unmarshal([]byte(line), &entry)
		if err != nil {
			log.Println("Failed to unmarshal:", err)
			continue
		}

		PushLogToDB(&entry)

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
	}
}

func PushLogToDB(entry *models.LogEntry) {
	if entry == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var project models.Project
	err := db.ProjectCollection.FindOne(ctx, bson.M{"api_key": entry.ApiKey}).Decode(&project)
	fmt.Println("Project Data:", project.UserID)
	if err != nil {
		return
	}

	var existingLog models.Log
	err = db.LogCollection.FindOne(ctx, bson.M{
		"message":    entry.Message,
		"level":      entry.Level,
		"timestamp":  entry.Timestamp,
		"project_id": project.ProjectID,
	}).Decode(&existingLog)

	if err == nil {
		fmt.Println("Duplicate log detected, skipping. Message:", entry.Message[:50])
		return
	}

	logEntry := models.Log{
		LogID:     bson.NewObjectID().Hex(),
		Level:     entry.Level,
		Message:   entry.Message,
		Service:   entry.Service,
		TimeStamp: entry.Timestamp,
		ProjectID: project.ProjectID,
		UserID:    project.UserID,
		Runtime: models.Runtime{
			File: entry.Runtime.File,
			Line: entry.Runtime.Line,
			Fn:   entry.Runtime.Fn,
		},
		Project:   project,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	fmt.Print("Log: ", logEntry)

	result, err := db.LogCollection.InsertOne(ctx, logEntry)
	if err != nil {
		fmt.Println("DB insert error:", err)
		return
	}

	if result.Acknowledged {
		db.ProjectCollection.UpdateOne(ctx, bson.M{"project_id": project.ProjectID}, bson.M{"$inc": bson.M{"log_counts": 1}})
		var user models.User

		fmt.Println("userId:", logEntry.UserID)
		err = db.UserCollection.FindOne(ctx, bson.M{"user_id": project.UserID}).Decode(&user)
		if err != nil {
			fmt.Println("Fetching project err: ", err)
			return
		}
		fmt.Println("User:", user)

		workers.SendEmailQueue(user.Email, user.UserID, logEntry.LogID)
	}

	fmt.Println("saved log:", result)

}
