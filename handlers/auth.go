package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/shareed2k/goth_fiber"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/service"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AuthRouteHandler struct {
	log           hclog.Logger
	userDataStore data.UserDataStore
	jwtManager    *service.JwtManager
}

func NewAuthRouteHandler(l hclog.Logger, userDataStore data.UserDataStore, jwtManager *service.JwtManager) *AuthRouteHandler {
	return &AuthRouteHandler{
		log:           l,
		userDataStore: userDataStore,
		jwtManager:    jwtManager,
	}
}

func (a AuthRouteHandler) OauthProviderCallback(c *fiber.Ctx) error {
	gothUser, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		a.log.Error("unable to complete auth\n" + err.Error())
		return c.Redirect("http://localhost:5050")
	}

	// Check existing user
	user, err := a.userDataStore.FindByProviderAndEmail(gothUser.Provider, gothUser.Email)
	if err != nil {
		a.log.Error("unable to find existing user\n" + err.Error())
		return c.Redirect("http://localhost:5050")
	}

	if user == nil {
		// Create new user
		username := a.getUserName(gothUser.NickName, gothUser.FirstName, gothUser.Name, gothUser.Email)
		user, err = a.userDataStore.Create(gothUser.Email, username, gothUser.Provider, gothUser.AvatarURL)
		if err != nil {
			a.log.Error("unable to create user\n" + err.Error())
			return c.Redirect("http://localhost:5050")
		}
	}

	// Create JWT
	token, err := a.jwtManager.GenerateUserJWT(strconv.Itoa(int(user.ID)))
	if err != nil {
		a.log.Error("unable to generate JWT\n" + err.Error())
		return c.Redirect("http://localhost:5050")
	}

	// Create cookie and attach
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Secure:   false,
		HTTPOnly: true,
		SameSite: "lax",
	})

	// Goth Session
	//if err := goth_fiber.StoreInSession("userId", strconv.Itoa(int(user.ID)), c); err != nil {
	//	a.log.Error("unable to store userId in session\n" + err.Error())
	//	return c.Redirect("http://localhost:5050")
	//}

	return c.Redirect("http://localhost:5050")
}

func (a *AuthRouteHandler) GetUserInfo(c *fiber.Ctx) error {
	token := c.Cookies("token", "")
	if len(token) == 0 {
		return NewHTTPError(a.log, fiber.StatusUnauthorized, "Token is not present", fmt.Errorf("no token is found"))
	}

	userIdString, err := a.jwtManager.ValidateUserJWT(token)
	if err != nil {
		c.ClearCookie("token")
		a.log.Error(err.Error())
		return NewHTTPError(a.log, fiber.StatusUnauthorized, "Invalid token", err)
	}

	// Get user
	userIdInt, err := strconv.Atoi(userIdString)
	if err != nil {
		c.ClearCookie("token")
		return NewHTTPError(a.log, fiber.StatusInternalServerError, "Unable to convert userId string to int", err)
	}
	user, err := a.userDataStore.FindById(uint(userIdInt))
	if err != nil {
		c.ClearCookie("token")
		return NewHTTPError(a.log, fiber.StatusInternalServerError, "Unable to find user in db", err)
	} else if user == nil {
		c.ClearCookie("token")
		return NewHTTPError(a.log, fiber.StatusUnauthorized, "User id not found", err)
	}

	return c.JSON(user)
}

func (a *AuthRouteHandler) Logout(c *fiber.Ctx) error {
	_ = goth_fiber.Logout(c)
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Expires:  time.Now().Add(-(time.Hour * 24)),
		Secure:   false,
		HTTPOnly: true,
		SameSite: "lax",
	})
	return c.JSON(fiber.Map{
		"success": true,
	})
}

func (a *AuthRouteHandler) getUserName(nickname, firstname, name, email string) string {
	if len(nickname) != 0 {
		return nickname
	}
	if len(firstname) != 0 {
		return firstname
	}
	if len(name) != 0 {
		return strings.Split(name, " ")[0]
	}

	emailRegex := regexp.MustCompile("(.+)@.+")
	subMatch := emailRegex.FindStringSubmatch(email)
	return subMatch[1]
}
