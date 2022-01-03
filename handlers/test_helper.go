package handlers

import (
	"github.com/stretchr/testify/mock"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"io"
)

type MockImageDataStore struct {
	mock.Mock
}

func (m *MockImageDataStore) Create(image *model.Image) error {
	args := m.Called(image)
	return args.Error(0)
}

func (m *MockImageDataStore) FindByUserID(userID uint) (*[]model.Image, error) {
	args := m.Called(userID)
	return args.Get(0).(*[]model.Image), args.Error(1)
}

func (m *MockImageDataStore) FindByID(imageID uint) (*model.Image, error) {
	args := m.Called(imageID)
	return args.Get(0).(*model.Image), args.Error(1)
}

func (m *MockImageDataStore) DeleteByID(imageId uint) error {
	args := m.Called(imageId)
	return args.Error(0)
}

type MockImageStorageManager struct {
	mock.Mock
}

func (m *MockImageStorageManager) UploadImage(fileName string, file io.ReadSeekCloser) error {
	args := m.Called(fileName, file)
	return args.Error(0)
}

func (m *MockImageStorageManager) DeleteImage(fileName string) error {
	args := m.Called(fileName)
	return args.Error(0)
}

type MockJWTManager struct {
	mock.Mock
}

func (m *MockJWTManager) Generate(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) Validate(tokenString string) (string, error) {
	args := m.Called(tokenString)
	return args.String(0), args.Error(1)
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) GenerateFileToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) GenerateFileID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) GenerateImageToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) GenerateAPIToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
