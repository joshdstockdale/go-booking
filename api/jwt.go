package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joshdstockdale/go-booking/db"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Println("Token not present in the header")
			return ErrUnAuthorized()
		}
		claims, err := validateToken(token[0])
		if err != nil {
			fmt.Println("Token not valid", err)
			return ErrUnAuthorized()
		}
		expires := claims["expires"]
		timeTime, err := time.Parse(time.RFC3339, expires.(string))
		if err != nil {
			return err
		}
		if time.Now().After(timeTime) {
			return NewError(http.StatusUnauthorized, "Token expired.")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrUnAuthorized()
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnAuthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Failed to parse JWT Token:", err)
		return nil, ErrUnAuthorized()
	}
	if !token.Valid {
		fmt.Println("Invalid Token")
		return nil, ErrUnAuthorized()
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, ErrUnAuthorized()
}
