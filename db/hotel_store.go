package db

import (
	"context"

	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const hotelsColl = "hotels"

type HotelQueryParams struct {
	Pagination
	Rating   int
	Location string
}

type HotelStore interface {
	Dropper

	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	AddRoom(context.Context, *types.Room) error
	List(context.Context, HotelQueryParams) ([]*types.Hotel, error)
	GetByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(c *mongo.Client, dbName string) *MongoHotelStore {
	coll := c.Database(dbName).Collection(hotelsColl)

	model := mongo.IndexModel{Keys: bson.M{"location": "text"}}
	_, err := coll.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		panic(err)
	}

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

func (s *MongoHotelStore) List(ctx context.Context, queryParams HotelQueryParams) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip((queryParams.Page - 1) * queryParams.Limit)
	opts.SetLimit(queryParams.Limit)

	filter := bson.M{}
	if queryParams.Rating != 0 {
		filter["rating"] = queryParams.Rating
	}
	if queryParams.Location != "" {
		filter["$text"] = bson.M{
			"$search": queryParams.Location,
		}
	}

	res, err := s.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err := res.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) GetByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := ToObjectID(id)
	if err != nil {
		return nil, err
	}

	var hotel types.Hotel
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel)
	if err != nil {
		return nil, err
	}

	return &hotel, nil
}
