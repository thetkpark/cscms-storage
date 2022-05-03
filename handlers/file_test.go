package handlers

import (
	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"go.uber.org/zap"
	"mime/multipart"
	"testing"
	"time"
)

type FileHandlerTestSuite struct {
	suite.Suite
	handler          *FileRoutesHandler
	store            *MockFileDataStore
	storage          *MockDiskStorageManager
	encrypt          *MockEncryptionManager
	token            *MockTokenManager
	ctx              *MockContext
	maxStoreDuration time.Duration
}

func TestFileRoutesHandler(t *testing.T) {
	suite.Run(t, new(FileHandlerTestSuite))
}

func (s *FileHandlerTestSuite) SetupTest() {
	s.store = new(MockFileDataStore)
	s.storage = new(MockDiskStorageManager)
	s.encrypt = new(MockEncryptionManager)
	s.token = new(MockTokenManager)
	s.maxStoreDuration = time.Hour
	s.handler = NewFileRoutesHandler(zap.NewNop().Sugar(), s.encrypt, s.store, s.storage, s.token, s.maxStoreDuration)
	s.ctx = new(MockContext)
}

func (s *FileHandlerTestSuite) TestAnonymousNoSlugFileUpload() {
	fileInfo := &model.File{
		ID:        faker.Password(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiredAt: time.Now().Add(s.maxStoreDuration),
		Token:     faker.Password(),
		Filename:  faker.Username(),
		FileSize:  20 << 20,
		Visited:   0,
		FileType:  "application/pdf",
		Encrypted: false,
	}
	multipartFile := &multipart.FileHeader{
		Filename: fileInfo.Filename,
		Size:     int64(fileInfo.FileSize),
	}
	file, err := multipartFile.Open()
	require.NoError(s.T(), err)
	s.ctx.On("FormFile", "file").Return(multipartFile)
	s.token.On("GenerateFileToken").Return(fileInfo.Token)
	s.ctx.On("Query", "slug", fileInfo.Token).Return(fileInfo.Token)
	s.store.On("FindByToken", fileInfo.Token).Return(nil, nil)
	s.ctx.On("Query", "duration", "").Return("")
	s.token.On("GenerateFileID").Return(fileInfo.ID)
	s.ctx.On("UserContext").Return(nil)
	s.storage.On("WriteToNewFile", fileInfo.ID, file).Return(nil)
	s.store.On("Create", fileInfo).Return(nil)
	s.ctx.On("Status", fiber.StatusCreated).Return(s.ctx)
	s.ctx.On("JSON", fileInfo).Return(nil)

	require.NoError(s.T(), s.handler.UploadFile(s.ctx))
}
