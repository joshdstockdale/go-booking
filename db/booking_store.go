package db

import (
	"context"
	"os"

	"github.com/joshdstockdale/go-booking/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	Insert(context.Context, *types.Booking) (*types.Booking, error)
	Get(context.Context, Map) ([]*types.Booking, error)
	GetByID(context.Context, string) (*types.Booking, error)
	Update(context.Context, string, Map) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	dbName := os.Getenv(MongoDBEnvName)
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbName).Collection("bookings"),
	}
}

func (s *MongoBookingStore) Update(ctx context.Context, id string, update Map) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	m := bson.M{"$set": update}
	_, err = s.coll.UpdateByID(ctx, oid, m)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoBookingStore) GetByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *MongoBookingStore) Get(ctx context.Context, filter Map) ([]*types.Booking, error) {
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

func (s *MongoBookingStore) Insert(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = resp.InsertedID.(primitive.ObjectID)
	return booking, nil

}
