package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf("mongodb+srv://%s:%s@thirfty.kzge54i.mongodb.net/?retryWrites=true&w=majority&appName=THIRFTY",
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
		),
	))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb")

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)

	return collection
}
