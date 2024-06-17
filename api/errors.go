package api

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
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

func ErrNotResourceNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  res + " resource not found",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "Invalid JSON Request",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "Invalid Id Given",
	}
}

func ErrUnAuthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized request",
	}
}
