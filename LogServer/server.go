package logserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/EsanSamuel/sensory/db"
	"github.com/EsanSamuel/sensory/jobs/workers"
	"github.com/EsanSamuel/sensory/models"
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

func Initialize_Log() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// WebSocket endpoint
	http.HandleFunc("/logs", handleWebSocket)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("WebSocket Log Server listening on port:", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Accepted connection from:", r.RemoteAddr)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var entry models.LogEntry
		err = json.Unmarshal(message, &entry)
		if err != nil {
			log.Println("Failed to unmarshal:", err)
			continue
		}

		PushLogToDB(&entry)
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
