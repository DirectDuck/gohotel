package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi      = "mongodb://admin:admin@localhost:27017"
	dbName     = "hotel-reservation"
	testdbName = "hotel-reservation-test"
)

func GetDatabaseClient() *mongo.Client {
	dbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	return dbClient
}

func GetDatabase() *mongo.Database {
	return GetDatabaseClient().Database(dbName)
}

func GetTestDatabase() *mongo.Database {
	return GetDatabaseClient().Database(testdbName)
}
