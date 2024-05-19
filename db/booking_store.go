package db

import (
	"context"

	"github.com/hhanri/ghotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	Insert(context.Context, *types.Booking) (*types.Booking, error)
	GetAll(context.Context, interface{}) ([]*types.Booking, error)
	GetByID(context.Context, string) (*types.Booking, error)
	Cancel(context.Context, string) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(c *mongo.Client, dbName string) *MongoBookingStore {
	coll := c.Database(dbName).Collection(bookingColl)
	return &MongoBookingStore{
		client: c,
		coll:   coll,
	}
}

func (s *MongoBookingStore) Insert(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = ObjectIdToString(resp.InsertedID)
	return booking, nil
}

func (s *MongoBookingStore) GetAll(ctx context.Context, filter interface{}) ([]*types.Booking, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := resp.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *MongoBookingStore) GetByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := ToObjectID(id)
	if err != nil {
		return nil, err
	}

	var booking types.Booking
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *MongoBookingStore) Cancel(ctx context.Context, id string) error {
	oid, err := ToObjectID(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"canceled": true,
		},
	}

	_, err = s.coll.UpdateByID(ctx, oid, update)
	return err
}

func (s *MongoBookingStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}
