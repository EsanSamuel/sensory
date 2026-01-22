package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env file not found!")
	}
	mongodb_uri := os.Getenv("mongodb_uri")

	if mongodb_uri == "" {
		fmt.Println("Mongodb uri is empty")
		return nil
	}

	clientOptions := options.Client().ApplyURI(mongodb_uri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("MongoDB is connected!")

	return client
}

var client *mongo.Client = ConnectDB()

func CollectionName(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env file not found!")
	}
	database_name := os.Getenv("database_name")

	if database_name == "" {
		fmt.Println("Database name is empty")
		return nil
	}

	collection := client.Database(database_name).Collection(collectionName)

	if collection == nil {
		return nil
	}

	return collection
}

var UserCollection = CollectionName("users")
var ProjectCollection = CollectionName("projects")
var LogCollection = CollectionName("logs")
