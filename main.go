package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joshdstockdale/go-booking/api"
	"github.com/joshdstockdale/go-booking/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userColl = "users"

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	coll := client.Database(dbname).Collection(userColl)
	user := types.User{
		FirstName: "Josh",
		LastName:  "Stock",
	}
	res, err := coll.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	listenAddr := flag.String("listenAddr", ":5000", "Listen address of the API Server")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)
	app.Listen(*listenAddr)
}
