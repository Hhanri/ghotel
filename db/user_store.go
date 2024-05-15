package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersColl = "users"

type UserStore interface {
	Dropper

	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, types.UpdateUserParams) error
	GetUserByEmail(context.Context, string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(c *mongo.Client, dbName string) *MongoUserStore {
	coll := c.Database(dbName).Collection(usersColl)
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

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("not found")
		}

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
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = ObjectIdToString(res.InsertedID)
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := ToObjectID(id)
	if err != nil {
		return err
	}

	// TODO: maybe handle if we didn't delete any user
	_, err = s.coll.DeleteOne(
		ctx,
		bson.M{"_id": oid},
	)
	return err
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, params types.UpdateUserParams) error {
	oid, err := ToObjectID(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": params.ToBSON(),
	}

	_, err = s.coll.UpdateByID(ctx, oid, update)
	return err
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("not found")
		}

		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping user collection ---")
	return s.coll.Drop(ctx)
}
