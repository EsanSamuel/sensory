package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	logserver "github.com/EsanSamuel/sensory/LogServer"
	"github.com/EsanSamuel/sensory/controllers"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env file not found!")
	}

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to sensory api"})
	})

	r.POST("/register", controllers.RegisterUser())
	r.GET("/user/:userId", controllers.GetUser())
	r.POST("project", controllers.CreateProject())

	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatal("Error starting server", err)
		}
	}()

	logserver.Initialize_Log()
}
