package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DBNAME = "ghotel"
const TestDBNAME = "ghotel-test"

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

func ToObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func ObjectIdToString(id interface{}) string {
	return id.(primitive.ObjectID).Hex()
}

type Dropper interface {
	Drop(context.Context) error
}

func NewMongoClient(dbUri string) (*mongo.Client, error) {
	return mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(dbUri),
	)
}
