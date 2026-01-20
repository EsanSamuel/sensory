package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/EsanSamuel/sensory/db"
	"github.com/EsanSamuel/sensory/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GenerateApiKey() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Error generating key")
		return ""
	}

	return hex.EncodeToString(b)
}

func CreateProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var project models.Project
		validate := validator.New()

		if err := c.ShouldBindJSON(&project); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			log.Println(err)
			return
		}

		if err := validate.Struct(project); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			log.Println(err)
			return
		}

		project.ProjectID = bson.NewObjectID().Hex()
		project.CreatedAt = time.Now()
		project.UpdatedAt = time.Now()

		result, err := db.ProjectCollection.InsertOne(ctx, project)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error creating project",
				"details": err.Error(),
			})
			log.Println(err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"project": result})
	}
}

func GenerateProjectApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		projectId := c.Param("projectId")
		api_key := GenerateApiKey()

		updateApiKey := bson.M{
			"$set": bson.M{
				"api_key":    api_key,
				"updated_at": time.Now(),
			},
		}

		result, err := db.ProjectCollection.UpdateOne(ctx, bson.M{"project_id": projectId}, updateApiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating apikey", "message": err})
			log.Println(err)
			return
		}
		if result.Acknowledged {
			c.JSON(http.StatusCreated, gin.H{"apikey": api_key})
		}

	}
}

func GetProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("userId")
		var projects []models.Project

		cursor, err := db.ProjectCollection.Find(ctx, bson.M{"user_id": userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fetching products failed", "details": err.Error()})
		}

		if err := cursor.All(ctx, &projects); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "decoding products failed", "details": err.Error()})
		}

		c.JSON(http.StatusOK, projects)

	}
}
