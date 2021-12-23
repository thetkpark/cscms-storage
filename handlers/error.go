package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewHTTPError(log *zap.SugaredLogger, code int, message string, error error) error {
	if code == fiber.StatusInternalServerError {
		log.Errorw(message, "error", error.Error())
	}
	return fiber.NewError(code, message)
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
