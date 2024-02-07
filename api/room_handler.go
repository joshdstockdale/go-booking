package api

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joshdstockdale/go-booking/db"
	"github.com/joshdstockdale/go-booking/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate  time.Time `json:"fromDate"`
	ToDate    time.Time `json:"toDate"`
	NumGuests int       `json:"numGuests"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
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
	booking := types.Booking{
		UserID:    user.ID,
		RoomID:    oid,
		FromDate:  params.FromDate,
		ToDate:    params.ToDate,
		NumGuests: params.NumGuests,
	}
	return c.JSON(booking)
}
