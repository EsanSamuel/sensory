package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EsanSamuel/sensory/db"
	"github.com/EsanSamuel/sensory/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Will Use clerk auth
func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		countUsers, err := db.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error ": "error counting document", "details": err})
		}

		if countUsers > 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "User already exists!", "error": err})
		}

		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		result, err := db.UserCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user", "message": err})
		}

		c.JSON(http.StatusCreated, gin.H{"user": result})
	}

}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("userId")

		var user models.User

		err := db.UserCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error finding user", "detail": err})
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
