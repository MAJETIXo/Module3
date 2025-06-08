package dto

import (
	"github.com/gofiber/fiber/v2"
)

const (
	FieldBadFormat     = "BAD_FORMAT"
	FieldIncorrect     = "INCORRECT"
	NotFound           = "NOT_FOUND"
	ServiceUnavailable = "UNAVAILABLE"
	InternalError      = "Service is currently unavailable. Please try again later."
)

type Response struct {
	Status string `json:"status"`
	Error  *Error `json:"error"`
	Data   any    `json:"data"`
}

type Error struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

func BadResponseError(ctx *fiber.Ctx, code, desk string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: code,
			Desc: desk,
		},
	})
}

func IternalServerError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: ServiceUnavailable,
			Desc: InternalError,
		},
	})
}

func NotFoundError(ctx *fiber.Ctx, desc string) error {
	return ctx.Status(fiber.StatusNotFound).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: NotFound,
			Desc: desc,
		},
	})
}
