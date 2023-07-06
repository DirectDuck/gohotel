package db

import (
	"context"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoUserColl     = "users"
	mongoHotelsColl   = "hotels"
	mongoRoomsColl    = "rooms"
	mongoBookingsColl = "bookings"
)

func GetMongoDBClient() *mongo.Client {
	options := options.Client().ApplyURI(os.Getenv("MONGO_DB_URI"))
	if strings.ToLower(os.Getenv("MONGO_DB_LOG_QUERIES")) == "true" {
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				log.Print(evt.Command)
			},
		}
		options = options.SetMonitor(cmdMonitor)
	}
	dbClient, err := mongo.Connect(context.Background(), options)
	if err != nil {
		log.Fatal(err)
	}
	return dbClient
}

type DB struct {
	mongoDBConn *mongo.Database
	Users       *MongoStore
	Hotels      *MongoStore
	Rooms       *MongoStore
	Bookings    *MongoStore
}

func GetDatabase() *DB {
	mongoDB := GetMongoDBClient().Database(os.Getenv("MONGO_DB_NAME"))
	db := &DB{
		mongoDBConn: mongoDB,
		Users:       &MongoStore{Coll: mongoDB.Collection(mongoUserColl)},
		Hotels:      &MongoStore{Coll: mongoDB.Collection(mongoHotelsColl)},
		Rooms:       &MongoStore{Coll: mongoDB.Collection(mongoRoomsColl)},
		Bookings:    &MongoStore{Coll: mongoDB.Collection(mongoBookingsColl)},
	}
	return db
}

func GetTestDatabase() *DB {
	mongoDB := GetMongoDBClient().Database(os.Getenv("MONGO_DB_TEST_NAME"))
	db := &DB{
		mongoDBConn: mongoDB,
		Users:       &MongoStore{Coll: mongoDB.Collection(mongoUserColl)},
		Hotels:      &MongoStore{Coll: mongoDB.Collection(mongoHotelsColl)},
		Rooms:       &MongoStore{Coll: mongoDB.Collection(mongoRoomsColl)},
		Bookings:    &MongoStore{Coll: mongoDB.Collection(mongoBookingsColl)},
	}
	return db
}

func (self *DB) Drop(ctx context.Context) error {
	return self.mongoDBConn.Drop(ctx)
}
