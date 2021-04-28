package config

import (
	"context"
	"log"
	"time"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "github.com/gofiber/fiber/v2"
)

// MongoInstance contains the Mongo client and database objects
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

// Database settings (insert your own database name and connection URI)
const dbName = "golang-cowhite-blog"
const mongoURI = "mongodb://localhost:27017/" + dbName


// Connect configures the MongoDB client and initializes the database connection.
// Source: https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup
func Connect() MongoInstance {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		log.Fatal(err)
	}

	mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	return mg;
}