package api

import (
	"github.com/joshdstockdale/go-booking/types"

	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Josh",
		LastName:  "Stockdale",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("me")
}
