package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return ctx.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return ctx.Status(apiError.Code).JSON(apiError)

}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}
func ErrUnAuthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "Unathorized",
	}
}
func ErrNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  res + " Not Found",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "Invalid JSON request",
	}
}
func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "Invalid id given",
	}
}
