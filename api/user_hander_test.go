package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joshdstockdale/go-booking/db"
	"github.com/joshdstockdale/go-booking/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testdburi = "mongodb://localhost:27017"
	dbname    = "hotel-reservation-test"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.InsertUserParams{
		Email:     "test@test.com",
		FirstName: "Testing",
		LastName:  "Testerson",
		Password:  "asdf1234",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(res.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("user.ID is not set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Encrypted Password should not be returned")
	}
	assertUserProps(t, "firstName", user.FirstName, params.FirstName)
	assertUserProps(t, "lastName", user.LastName, params.LastName)
	assertUserProps(t, "email", user.Email, params.Email)
}

func assertUserProps(t testing.TB, prop string, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("For prop (%v), got %v but wanted %v", prop, got, want)
	}
}
