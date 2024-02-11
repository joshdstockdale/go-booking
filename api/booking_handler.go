package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joshdstockdale/go-booking/db"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetByID(c.Context(), id)
	if err != nil {
		return ErrNotFound("Booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnAuthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnAuthorized()
	}
	if err := h.store.Booking.Update(c.Context(), id, bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResponse{Type: "success", Msg: "Updated"})
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.Get(c.Context(), bson.M{})
	if err != nil {
		return ErrNotFound("Bookings")
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetByID(c.Context(), id)
	if err != nil {
		return ErrNotFound("Booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnAuthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnAuthorized()
	}

	return c.JSON(booking)
}
