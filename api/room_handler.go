package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joshdstockdale/go-booking/db"
	"github.com/joshdstockdale/go-booking/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate  time.Time `json:"fromDate"`
	ToDate    time.Time `json:"toDate"`
	NumGuests int       `json:"numGuests"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.ToDate) {
		return fmt.Errorf("Cannot book a room in the past.")
	}
	if p.FromDate.After(p.ToDate) {
		return fmt.Errorf("To Date must be before From Date")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.Get(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}

	roomID := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type: "error",
			Msg:  "Internal Server Error",
		})
	}
	ok, err = h.isRoomAvailable(c.Context(), oid, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(
			genericResponse{
				Type: "error",
				Msg:  fmt.Sprintf("Room is already booked for those dates."),
			},
		)
	}
	booking := types.Booking{
		UserID:    user.ID,
		RoomID:    oid,
		FromDate:  params.FromDate,
		ToDate:    params.ToDate,
		NumGuests: params.NumGuests,
	}
	inserted, err := h.store.Booking.Insert(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailable(c context.Context, oid primitive.ObjectID, params BookRoomParams) (bool, error) {

	where := bson.M{
		"roomID": oid,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"toDate": bson.M{
			"$lte": params.ToDate,
		},
	}
	bookings, err := h.store.Booking.Get(c, where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}
