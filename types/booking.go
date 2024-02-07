package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID    primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NumGuests int                `bson:"numGuests" json:"numGuests"`
	FromDate  time.Time          `bson:"from" json:"from"`
	ToDate    time.Time          `bson:"toDate" json:"toDate"`
}
