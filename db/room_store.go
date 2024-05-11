package db

import (
	"context"

	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomsColl = "rooms"

type RoomStore interface {
	Insert(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(c *mongo.Client, dbName string, hotelStore HotelStore) *MongoRoomStore {
	coll := c.Database(dbName).Collection(roomsColl)
	return &MongoRoomStore{
		client:     c,
		coll:       coll,
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) Insert(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = ObjectIdToString(res.InsertedID)

	err = s.HotelStore.AddRoom(ctx, room)
	if err != nil {
		return nil, err
	}

	return room, nil
}