package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joshdstockdale/go-booking/db"
	"github.com/joshdstockdale/go-booking/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}
type genericResponse struct {
	Type string
	Msg  string
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResponse{
		Type: "error",
		Msg:  "Invalid Credentials",
	})
}

// A handler should only do:
// - serialization of incoming request
// - fetch data
// - call some business logic (or 3rd party lib)k
// - return data
func (h *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}
	user, err := h.userStore.GetByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsPasswordValid(user.EncryptedPassword, authParams.Password) {
		return invalidCredentials(c)
	}
	token, err := types.CreateTokenFromUser(user)
	if err != nil {
		return err
	}
	resp := AuthResponse{
		User:  user,
		Token: token,
	}
	return c.JSON(resp)
}
