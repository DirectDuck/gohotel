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

type Store struct {
	Users  UserStore
	Hotels HotelStore
	Rooms  RoomStore
}

type MongoDB struct {
	*mongo.Database
	Store *Store
}

func NewStore(dbSrc *MongoDB) *Store {
	store := &Store{
		Users:  NewMongoUserStore(dbSrc),
		Hotels: NewMongoHotelStore(dbSrc),
		Rooms:  NewMongoRoomStore(dbSrc),
	}
	return store
}

func GetDatabase() *MongoDB {
	baseMongo := &MongoDB{
		Database: GetDatabaseClient().Database(dbName),
	}
	baseMongo.Store = NewStore(baseMongo)
	return baseMongo
}

func GetTestDatabase() *MongoDB {
	baseMongo := &MongoDB{
		Database: GetDatabaseClient().Database(testdbName),
	}
	baseMongo.Store = NewStore(baseMongo)
	return baseMongo
}
