package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joshdstockdale/go-booking/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}
	filter := db.Map{"hotelID": oid}
	rooms, err := h.store.Room.Get(c.Context(), filter)
	if err != nil {
		return ErrNotFound("Rooms")
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetByID(c.Context(), id)
	if err != nil {
		return ErrNotFound("Hotel")
	}
	return c.JSON(hotel)
}

type ResourceResp struct {
	Total int
	Page  int
	Data  any
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	filter := db.Map{
		"rating": params.Rating,
	}
	hotels, err := h.store.Hotel.Get(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrNotFound("Hotels")
	}
	resp := ResourceResp{
		Data:  hotels,
		Page:  int(params.Pagination.Page),
		Total: len(hotels),
	}
	return c.JSON(resp)
}
