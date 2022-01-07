package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thetkpark/cscms-temp-storage/handlers"
	"go.uber.org/zap"
)

func NewFiberRouter() *fiber.App {
	return fiber.New(fiber.Config{
		BodyLimit: 150 << 20,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default to 500
			code := fiber.StatusInternalServerError
			message := err.Error()

			// Check if error is fiber.Error type
			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
				message = e.Message
			}

			return c.Status(code).JSON(fiber.Map{
				"code":    code,
				"message": message,
			})
		},
	})
}

type HandlerFunc func(c handlers.Context) error

func CreateFiberHandler(handlerFunc HandlerFunc, logger *zap.SugaredLogger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return handlerFunc(&FiberContext{
			logger: logger,
			Ctx:    c,
		})
	}
}

type FiberContext struct {
	*fiber.Ctx
	logger *zap.SugaredLogger
}

func (c *FiberContext) Status(code int) handlers.Context {
	c.Ctx.Status(code)
	return c
}

func (c *FiberContext) Redirect(location string) error {
	return c.Redirect(location)
}

func (c *FiberContext) Error(code int, message string, error error) error {
	if code == fiber.StatusInternalServerError {
		c.logger.Errorw(message, "error", error.Error())
	}
	return c.Ctx.Status(code).JSON(&handlers.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
