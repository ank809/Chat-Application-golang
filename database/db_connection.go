package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DbClient()

func DbClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := godotenv.Load(); err != nil {
		log.Println(err)
		return nil
	}
	url := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Println(err)
		return nil
	}
	return client
}

func OpenCollection(c *mongo.Client, collection_name string) *mongo.Collection {
	coll := c.Database("Chat-Application").Collection(collection_name)
	return coll
}
