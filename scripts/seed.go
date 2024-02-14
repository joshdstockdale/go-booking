package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/joshdstockdale/go-booking/db"
	"github.com/joshdstockdale/go-booking/db/fixtures"
	"github.com/joshdstockdale/go-booking/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatal(err)
	}
	log.Println(os.LookupEnv("MONGO_DB_URL"))
	log.Println(os.LookupEnv("MONGO_DB_NAME"))
	var (
		ctx         = context.Background()
		mongoUrl    = os.Getenv("MONGO_DB_URL")
		mongoDbName = os.Getenv("MONGO_DB_NAME")
	)
	log.Println("AFTER", mongoDbName)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Fatal(err)
	}
	//if err := client.Database(mongoDbName).Drop(ctx); err != nil {
	//	log.Fatal(err)
	//}
	hotelStore := db.NewMongoHotelStore(client)
	log.Println("AFTER", mongoDbName)
	store := db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   db.NewMongoHotelStore(client),
	}
	user := fixtures.InsertUser(&store, "Josh", "NoAdmin", false)
	token, _ := types.CreateTokenFromUser(user)
	fmt.Println("--NonAdmin:", token)
	admin := fixtures.InsertUser(&store, "Josh", "Admin", true)
	token, _ = types.CreateTokenFromUser(admin)
	fmt.Println("--Admin:", token)
	hotel := fixtures.InsertHotel(&store, "BnB", "Mountains", 5)
	room := fixtures.InsertRoom(&store, "large", false, 211.12, hotel.ID)
	booking := fixtures.InsertBooking(
		&store, user.ID, room.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2),
	)
	fmt.Println("--Booking:", booking.ID)
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Hotel %d", i)
		location := fmt.Sprintf("loc %d", i)
		fixtures.InsertHotel(&store, name, location, rand.Intn(5)+1)
	}
}
