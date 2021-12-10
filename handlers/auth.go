package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/shareed2k/goth_fiber"
	"github.com/thetkpark/cscms-temp-storage/data"
	"time"
)

type AuthRouteHandler struct {
	log           hclog.Logger
	userDataStore data.UserDataStore
}

func NewAuthRouteHandler(l hclog.Logger, userDataStore data.UserDataStore) *AuthRouteHandler {
	return &AuthRouteHandler{
		log:           l,
		userDataStore: userDataStore,
	}
}

func (a AuthRouteHandler) OauthProviderCallback(c *fiber.Ctx) error {
	if gothUser, err := goth_fiber.CompleteUserAuth(c); err == nil {
		cookie := new(fiber.Cookie)
		cookie.Name = "token"
		cookie.Value = "jwt"
		cookie.Expires = time.Now().Add(30 * 24 * time.Hour)
		c.Cookie(cookie)
		return c.JSON(gothUser)
	} else {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
