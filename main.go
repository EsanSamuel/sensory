package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	logserver "github.com/EsanSamuel/sensory/LogServer"
	"github.com/EsanSamuel/sensory/controllers"
	"github.com/EsanSamuel/sensory/jobs/workers"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env file not found!")
	}

	r := gin.Default()

	workers.EmailWorker()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://sensory-frontend.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
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

	r.GET("/ws/logs", logserver.HandleWebSocketLogs)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // fallback for local development
	}


	go func() {
		fmt.Println("HTTP Server starting on port:", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatal("Error starting server:", err)
		}
	}()

	// Wait for shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	fmt.Println("Shutting down gracefully...")
	workers.StopEmailWorker()
}
