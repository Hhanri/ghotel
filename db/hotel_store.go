package db

import (
	"context"

	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelsColl = "hotels"

type HotelStore interface {
	Dropper

	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	AddRoom(context.Context, *types.Room) error
	List(context.Context, interface{}) ([]*types.Hotel, error)
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

	update := bson.M{
		"$push": bson.M{
			"rooms": room.ID,
		},
	}

	_, err = s.coll.UpdateByID(ctx, hotelID, update)
	return err
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}

func (s *MongoHotelStore) List(ctx context.Context, filter interface{}) ([]*types.Hotel, error) {
	res, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err := res.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}
