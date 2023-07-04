package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	dbClient, err := mongo.Connect(
		context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_DB_URI")),
	)
	if err != nil {
		log.Fatal(err)
	}
	return dbClient
}

type Store interface {
	Create(context.Context, interface{}) (primitive.ObjectID, error)
	Get(context.Context, interface{}, interface{}) (interface{}, error)
	GetOne(context.Context, interface{}, interface{}) (interface{}, error)
	GetOneByID(context.Context, primitive.ObjectID, interface{}) (interface{}, error)
	UpdateByID(context.Context, primitive.ObjectID, interface{}) error
	DeleteByID(context.Context, primitive.ObjectID) error
}

type DB struct {
	mongoDBConn *mongo.Database
	Users       Store
	Hotels      Store
	Rooms       Store
	Bookings    Store
}

func GetDatabase() *DB {
	mongoDB := GetMongoDBClient().Database(os.Getenv("MONGO_DB_NAME"))
	mongo := &DB{
		mongoDBConn: mongoDB,
		Users:       &MongoStore{Coll: mongoDB.Collection(mongoUserColl)},
		Hotels:      &MongoStore{Coll: mongoDB.Collection(mongoHotelsColl)},
		Rooms:       &MongoStore{Coll: mongoDB.Collection(mongoRoomsColl)},
		Bookings:    &MongoStore{Coll: mongoDB.Collection(mongoBookingsColl)},
	}
	return mongo
}

func GetTestDatabase() *DB {
	mongoDB := GetMongoDBClient().Database(os.Getenv("MONGO_TEST_DB_NAME"))
	mongo := &DB{
		mongoDBConn: mongoDB,
		Users:       &MongoStore{Coll: mongoDB.Collection(mongoUserColl)},
		Hotels:      &MongoStore{Coll: mongoDB.Collection(mongoHotelsColl)},
		Rooms:       &MongoStore{Coll: mongoDB.Collection(mongoRoomsColl)},
		Bookings:    &MongoStore{Coll: mongoDB.Collection(mongoBookingsColl)},
	}
	return mongo
}

func (self *DB) Drop(ctx context.Context) error {
	return self.mongoDBConn.Drop(ctx)
}
