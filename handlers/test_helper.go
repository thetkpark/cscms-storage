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

type MockFileDataStore struct {
	mock.Mock
}

func (m *MockFileDataStore) Create(file *model.File) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockFileDataStore) FindByID(fileID string) (*model.File, error) {
	args := m.Called(fileID)
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockFileDataStore) FindByToken(token string) (*model.File, error) {
	args := m.Called(token)
	return args.Get(0).(*model.File), args.Error(1)
}

func (m *MockFileDataStore) IncreaseVisited(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockFileDataStore) FindByUserID(userId uint) (*[]model.File, error) {
	args := m.Called(userId)
	return args.Get(0).(*[]model.File), args.Error(1)
}

func (m *MockFileDataStore) DeleteByID(fileId string) error {
	args := m.Called(fileId)
	return args.Error(0)
}

func (m *MockFileDataStore) UpdateToken(fileID string, newToken string) error {
	args := m.Called(fileID, newToken)
	return args.Error(0)
}

type MockDiskStorageManager struct {
	mock.Mock
}

func (m *MockDiskStorageManager) OpenFile(fileName string) (io.Reader, error) {
	args := m.Called(fileName)
	return args.Get(0).(io.Reader), args.Error(1)
}

func (m *MockDiskStorageManager) WriteToNewFile(fileName string, reader io.Reader) error {
	args := m.Called(fileName, reader)
	return args.Error(0)
}

func (m *MockDiskStorageManager) Exist(fileName string) (bool, error) {
	args := m.Called(fileName)
	return args.Bool(0), args.Error(1)
}

func (m *MockDiskStorageManager) ListFiles() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockDiskStorageManager) DeleteFile(fileName string) error {
	args := m.Called(fileName)
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

type MockEncryptionManager struct {
	mock.Mock
}

func (m *MockEncryptionManager) Encrypt(input io.Reader) (io.Reader, string, error) {
	args := m.Called(input)
	return args.Get(0).(io.Reader), args.String(1), args.Error(2)
}

func (m *MockEncryptionManager) Decrypt(input io.Reader, nonceString string) (io.Reader, error) {
	args := m.Called(input, nonceString)
	return args.Get(0).(io.Reader), args.Error(1)
}
