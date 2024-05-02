package db

import (
	"context"

	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersColl = "users"

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(c *mongo.Client) *MongoUserStore {
	coll := c.Database(DBNAME).Collection(usersColl)
	return &MongoUserStore{
		client: c,
		coll:   coll,
	}
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	var user types.User
	err := s.coll.FindOne(ctx, bson.M{"_id": ToObjectID(id)}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
