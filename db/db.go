package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const DBNAME = "ghotel"
const TestDBNAME = "ghotel-test"

func ToObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func ObjectIdToString(id interface{}) string {
	return id.(primitive.ObjectID).Hex()
}

type Dropper interface {
	Drop(context.Context) error
}
