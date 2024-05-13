package db

import (
	"context"

	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomsColl = "rooms"

type RoomStore interface {
	Insert(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, interface{}) ([]*types.Room, error)
	GetRoomsByID(context.Context, string) ([]*types.Room, error)
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

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter interface{}) ([]*types.Room, error) {
	res, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := res.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *MongoRoomStore) GetRoomsByID(ctx context.Context, id string) ([]*types.Room, error) {
	return s.GetRooms(ctx, bson.M{"hotelId": id})
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}
