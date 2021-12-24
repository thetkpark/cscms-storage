package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/service"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AuthRouteHandler struct {
	log           *zap.SugaredLogger
	userDataStore data.UserDataStore
	jwtManager    *service.JwtManager
	entrypoint    string
}

func NewAuthRouteHandler(l *zap.SugaredLogger, userDataStore data.UserDataStore, jwtManager *service.JwtManager, entry string) *AuthRouteHandler {
	return &AuthRouteHandler{
		log:           l,
		userDataStore: userDataStore,
		jwtManager:    jwtManager,
		entrypoint:    entry,
	}
}

func (a AuthRouteHandler) OauthProviderCallback(c *fiber.Ctx) error {
	gothUser, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		a.log.Error("unable to complete auth\n" + err.Error())
		return c.Redirect(a.entrypoint)
	}

	// Check existing user
	user, err := a.userDataStore.FindByProviderAndEmail(gothUser.Provider, gothUser.Email)
	if err != nil {
		a.log.Error("unable to find existing user\n" + err.Error())
		return c.Redirect(a.entrypoint)
	}

	if user == nil {
		// Create new user
		username := a.getUserName(gothUser.NickName, gothUser.FirstName, gothUser.Name, gothUser.Email)
		user, err = a.userDataStore.Create(gothUser.Email, username, gothUser.Provider, gothUser.AvatarURL)
		if err != nil {
			a.log.Error("unable to create user\n" + err.Error())
			return c.Redirect(a.entrypoint)
		}
	}

	// Create JWT
	token, err := a.jwtManager.GenerateUserJWT(strconv.Itoa(int(user.ID)))
	if err != nil {
		a.log.Error("unable to generate JWT\n" + err.Error())
		return c.Redirect(a.entrypoint)
	}

	// Create cookie and attach
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "lax",
	})

	// Goth Session
	//if err := goth_fiber.StoreInSession("userId", strconv.Itoa(int(user.ID)), c); err != nil {
	//	a.log.Error("unable to store userId in session\n" + err.Error())
	//	return c.Redirect("http://localhost:5050")
	//}

	return c.Redirect(a.entrypoint)
}

// GetUserInfo handlers
// @Summary Get user info
// @Description Get the user information
// @Tags Auth
// @Success  200 {object} model.User
// @Failure  401 {object}  handlers.ErrorResponse
// @Router /auth/user [get]
func (a *AuthRouteHandler) GetUserInfo(c *fiber.Ctx) error {
	user := c.UserContext().Value("user")
	if user == nil {
		return NewHTTPError(a.log, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	return c.JSON(user)
}

// Logout handlers
// @Summary Logout
// @Description Clear the cookie
// @Tags Auth
// @Success  200
// @Router /auth/logout [get]
func (a *AuthRouteHandler) Logout(c *fiber.Ctx) error {
	a.clearCookie(c)
	return c.SendStatus(fiber.StatusOK)
}

func (a *AuthRouteHandler) ParseUserFromCookie(c *fiber.Ctx) error {
	token := c.Cookies("token", "")
	if len(token) == 0 {
		return c.Next()
	}

	userIdString, err := a.jwtManager.ValidateUserJWT(token)
	if err != nil {
		a.clearCookie(c)
		return c.Next()
	}

	// Get user
	userIdInt, err := strconv.Atoi(userIdString)
	if err != nil {
		a.clearCookie(c)
		return NewHTTPError(a.log, fiber.StatusInternalServerError, "Unable to convert userId string to int", err)
	}
	user, err := a.userDataStore.FindById(uint(userIdInt))
	if err != nil || user == nil {
		a.clearCookie(c)
		return c.Next()
	}

	c.SetUserContext(context.WithValue(c.UserContext(), "user", user))
	return c.Next()
}

func (a *AuthRouteHandler) AuthenticatedOnly(c *fiber.Ctx) error {
	user := c.UserContext().Value("user")
	if user == nil {
		return NewHTTPError(a.log, fiber.StatusUnauthorized, "Unauthenticated", nil)
	}
	return c.Next()
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

func (a *AuthRouteHandler) clearCookie(c *fiber.Ctx) {
	_ = goth_fiber.Logout(c)
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Expires:  time.Now().Add(-(time.Hour * 24)),
		Secure:   false,
		HTTPOnly: true,
		SameSite: "lax",
	})
}
