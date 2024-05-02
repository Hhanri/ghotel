package db

import (
	"context"

	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersColl = "users"

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
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
	oid, err := ToObjectID(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err := cursor.Decode(&users); err != nil {
		return []*types.User{}, nil
	}

	return users, nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}
