package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const DBNAME = "ghotel"

func ToObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

type Dropper interface {
	Drop(context.Context) error
}
