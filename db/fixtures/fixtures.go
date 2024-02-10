package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joshdstockdale/go-booking/db"
	"github.com/joshdstockdale/go-booking/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertUser(store *db.Store, fname, lname string, isAdmin bool) *types.User {

	user, err := types.NewUserFromParams(types.InsertUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func InsertHotel(store *db.Store, name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
		Rooms:    []primitive.ObjectID{},
	}
	insertedHotel, err := store.Hotel.Insert(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func InsertRoom(store *db.Store, size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		Size:    size,
		Price:   price,
		Seaside: seaside,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.Insert(context.TODO(), &room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func InsertBooking(store *db.Store, userID, roomID primitive.ObjectID, numGuests int, from, to time.Time) *types.Booking {
	booking := types.Booking{
		NumGuests: numGuests,
		UserID:    userID,
		RoomID:    roomID,
		FromDate:  from,
		ToDate:    to,
	}
	insertedBooking, err := store.Booking.Insert(context.TODO(), &booking)

	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
