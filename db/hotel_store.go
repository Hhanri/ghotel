package db

import (
	"context"

	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelsColl = "hotels"

type HotelStore interface {
	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	AddRoom(context.Context, *types.Room) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(c *mongo.Client, dbName string) *MongoHotelStore {
	coll := c.Database(dbName).Collection(hotelsColl)
	return &MongoHotelStore{
		client: c,
		coll:   coll,
	}
}

func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = ObjectIdToString(res.InsertedID)
	return hotel, nil
}

func (s *MongoHotelStore) AddRoom(ctx context.Context, room *types.Room) error {
	hotelID, err := ToObjectID(room.HotelID)
	if err != nil {
		return err
	}

	roomID, err := ToObjectID(room.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{
			"rooms": roomID,
		},
	}

	_, err = s.coll.UpdateByID(ctx, hotelID, update)
	return err
}
