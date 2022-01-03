package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/router"
	"github.com/thetkpark/cscms-temp-storage/service/storage"
	"github.com/thetkpark/cscms-temp-storage/service/token"
	"go.uber.org/zap"
	"testing"
)

type ImageHandlerTestSuite struct {
	suite.Suite
	app     *fiber.App
	handler *ImageRouteHandler
	store   data.ImageDataStore
	storage storage.ImageManager
	token   token.Manager
}

func TestGormUserDataStore(t *testing.T) {
	suite.Run(t, new(ImageHandlerTestSuite))
}

func (s *ImageHandlerTestSuite) SetupTest() {
	s.app = router.NewFiberRouter()

	s.store = &MockImageDataStore{}
	s.storage = &MockImageStorageManager{}
	s.token = &MockTokenManager{}

	s.handler = &ImageRouteHandler{
		log:               zap.NewNop().Sugar(),
		imageDataStore:    s.store,
		imageStoreManager: s.storage,
		tokenManager:      s.token,
	}

	s.app.Post("/api/image", s.handler.UploadImage)
	s.app.Get("/api/image", s.handler.GetOwnImages)
	s.app.Delete("/api/image/:imageID", s.handler.DeleteImage)
	s.app.Patch("/api/image/own", s.handler.IsOwnImage, func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}

func (s *ImageHandlerTestSuite) AfterTest(_, _ string) {

}

func (s *ImageHandlerTestSuite) TestGetOwnImages() {
	//req := httptest.NewRequest(http.MethodGet, "")
}
