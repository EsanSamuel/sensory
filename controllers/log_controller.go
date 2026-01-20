package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EsanSamuel/sensory/db"
	"github.com/EsanSamuel/sensory/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("userId")

		pipeline := mongo.Pipeline{
			{{"$match", bson.D{
				{"user_id", userId},
			}}},
			{{"$lookup", bson.D{
				{"from", "projects"},
				{"localField", "project_id"},
				{"foreignField", "project_id"},
				{"as", "project"},
			}}},
			{{"$unwind", bson.D{
				{"path", "$project"},
				{"preserveNullAndEmptyArrays", true},
			}}},
		}

		cursor, err := db.LogCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "aggregation failed",
				"details": err.Error(),
			})
			return
		}

		var logs []models.Log
		if err := cursor.All(ctx, &logs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "decoding logs failed",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, logs)
	}
}

func GetLogsByProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		projectId := c.Param("projectId")

		pipeline := mongo.Pipeline{
			{{"$match", bson.D{
				{"project_id", projectId},
			}}},
			{{"$lookup", bson.D{
				{"from", "projects"},
				{"localField", "project_id"},
				{"foreignField", "project_id"},
				{"as", "project"},
			}}},
			{{"$unwind", bson.D{
				{"path", "$project"},
				{"preserveNullAndEmptyArrays", true},
			}}},
		}

		cursor, err := db.LogCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "aggregation failed",
				"details": err.Error(),
			})
			return
		}

		var logs []models.Log
		if err := cursor.All(ctx, &logs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "decoding logs failed",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, logs)
	}
}

func FilterLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("userId")
		levelQuery := c.Query("level")
		var logs []models.Log

		cursor, err := db.LogCollection.Find(ctx, bson.M{"user_id": userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fetching logs failed", "details": err.Error()})
		}

		if err := cursor.All(ctx, &logs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "decoding logs failed", "details": err.Error()})
		}

		level := make(map[string][]models.Log)

		for _, log := range logs {
			level[log.Level] = append(level[log.Level], log)
		}

		c.JSON(http.StatusOK, level[levelQuery])

	}
}

func GetLogById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var log models.Log
		logId := c.Param("logId")

		err := db.LogCollection.FindOne(ctx, bson.M{"log_id": logId}).Decode(&log)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fetching log failed", "details": err.Error()})
		}

		c.JSON(http.StatusOK, log)
	}
}
