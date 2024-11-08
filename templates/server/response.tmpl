package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Success(c *fiber.Ctx, data any, status ...int) error {
	code := fiber.StatusOK
	if len(status) > 0 {
		code = status[0]
	}
	return c.Status(code).JSON(fiber.Map{
		"code":   code,
		"status": "success",
		"data":   data,
	})
}

func Error(c *fiber.Ctx, err error, status ...int) error {
	code := fiber.StatusBadRequest
	if len(status) > 0 {
		code = status[0]
	}

	response := fiber.Map{
		"code":   code,
		"status": "error",
	}

	if apiResponseError, ok := err.(ApiResponseError); ok {
		response["error"] = apiResponseError.Parse()
	} else {
		response["error"] = (&ApiResponseError{MainError: err, Index: -1}).Parse()
	}

	return c.Status(code).JSON(response)
}
