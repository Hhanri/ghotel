package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const DBNAME = "ghotel"

func ToObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
