package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
)

func NewHTTPError(log hclog.Logger, code int, message string, error error) error {
	if code == 500 {
		log.Error(message, error.Error())
	}
	return fiber.NewError(code, message)
}
