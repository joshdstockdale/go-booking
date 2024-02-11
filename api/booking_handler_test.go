package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joshdstockdale/go-booking/db/fixtures"
	"github.com/joshdstockdale/go-booking/types"
)

func TestUserGetBooking(t *testing.T) {

	db := setup(t)
	defer db.teardown(t)
	var (
		nonAuthUser = fixtures.InsertUser(db.Store, "non", "auth", false)
		user        = fixtures.InsertUser(db.Store, "josh", "me", false)
		hotel       = fixtures.InsertHotel(db.Store, "Big Hotel", "a", 4)
		room        = fixtures.InsertRoom(db.Store, "small room", true, 22.3, hotel.ID)
		booking     = fixtures.InsertBooking(
			db.Store, user.ID, room.ID, 3, time.Now(), time.Now().AddDate(0, 0, 5),
		)
		app            = fiber.New()
		route          = app.Group("/", JWTAuthentication(db.User))
		BookingHandler = NewBookingHandler(db.Store)
	)
	route.Get("/:id", BookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	token, _ := types.CreateTokenFromUser(user)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected a 200 response but got %d", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	if bookingResp.ID != booking.ID {
		t.Fatalf("Expected  %s got %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.RoomID != booking.RoomID {
		t.Fatalf("Expected  %s got %s", booking.RoomID, bookingResp.RoomID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("Expected  %s got %s", booking.UserID, bookingResp.UserID)
	}

	// Non Authorized User
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	token, _ = types.CreateTokenFromUser(nonAuthUser)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 response but got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		adminUser = fixtures.InsertUser(db.Store, "admin", "admin", true)
		user      = fixtures.InsertUser(db.Store, "josh", "me", false)
		hotel     = fixtures.InsertHotel(db.Store, "Big Hotel", "a", 4)
		room      = fixtures.InsertRoom(db.Store, "small room", true, 22.3, hotel.ID)
		booking   = fixtures.InsertBooking(
			db.Store, user.ID, room.ID, 3, time.Now(), time.Now().AddDate(0, 0, 5),
		)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin          = app.Group("/", JWTAuthentication(db.User), AdminAuthentication)
		BookingHandler = NewBookingHandler(db.Store)
	)
	admin.Get("/", BookingHandler.HandleGetBookings)

	//Admin, should be able to access
	req := httptest.NewRequest("GET", "/", nil)
	token, _ := types.CreateTokenFromUser(adminUser)
	req.Header.Add("X-Api-Token", token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected a 200 response but got %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("Expected 1 booking got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("Expected  %s got %s", booking.ID, have.ID)
	}
	if have.RoomID != booking.RoomID {
		t.Fatalf("Expected  %s got %s", booking.RoomID, have.RoomID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("Expected  %s got %s", booking.UserID, have.UserID)
	}

	//Not an Admin, should not be able to access
	req = httptest.NewRequest("GET", "/", nil)
	token, _ = types.CreateTokenFromUser(user)
	req.Header.Add("X-Api-Token", token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected a Unauthorized response but got %d", resp.StatusCode)
	}
}
