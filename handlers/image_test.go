package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/service/storage"
	"github.com/thetkpark/cscms-temp-storage/service/token"
)

type ImageHandlerTestSuite struct {
	suite.Suite
	app     *fiber.App
	handler *ImageRouteHandler
	store   data.ImageDataStore
	storage storage.ImageManager
	token   token.Manager
}

//func TestGormUserDataStore(t *testing.T) {
//	suite.Run(t, new(ImageHandlerTestSuite))
//}

func (s *ImageHandlerTestSuite) SetupTest() {

}

func (s *ImageHandlerTestSuite) AfterTest(_, _ string) {

}

func (s *ImageHandlerTestSuite) TestGetOwnImages() {
	//req := httptest.NewRequest(http.MethodGet, "")
}
