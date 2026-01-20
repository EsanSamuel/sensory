package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to sensory api"})
	})

	r.POST("/register", controllers.RegisterUser())
	r.POST("/project", controllers.CreateProject())
	r.GET("/user/:userId", controllers.GetUser())
	r.GET("/projects/:userId", controllers.GetProjects())
	r.GET("/project/:projectId", controllers.GetProject())
	r.POST("/project/api_key/:projectId", controllers.GenerateProjectApiKey())
	r.GET("/logs/:userId", controllers.GetLogs())
	r.GET("/logs/level/:userId", controllers.FilterLogs())
	r.GET("/log/:logId", controllers.GetLogById())
	r.GET("/logs/project/:projectId", controllers.GetLogsByProject())

	go func() {
		if err := r.Run(":8000"); err != nil {
			log.Fatal("Error starting server", err)
		}
	}()

	logserver.Initialize_Log()
}
