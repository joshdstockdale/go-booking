package main

import (
	"context"
	"log"

	"github.com/joshdstockdale/go-booking/db"
	"github.com/joshdstockdale/go-booking/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(fname, lname, email string) {
	user, err := types.NewUserFromParams(types.InsertUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  "asdf1234",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name string, location string, rating int) {

	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
		Rooms:    []primitive.ObjectID{},
	}
	rooms := []types.Room{
		{
			Size:  "small",
			Price: 88.90,
		},
		{
			Size:  "normal",
			Price: 188.90,
		},
		{
			Size:  "large",
			Price: 298.90,
		},
	}
	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.Insert(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("Bed & Breakfast", "Grandfather Mt, NC", 5)
	seedHotel("Highrise Hotel", "Atlanta, GA", 4)
	seedHotel("Swanky Hotel", "New York, NY", 3)
	seedHotel("Plaid Hotel", "Seattle, WA", 2)
	seedUser("Josh", "Me", "josh@me.com")
}

func init() {
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
