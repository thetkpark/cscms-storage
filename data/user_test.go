package data

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"testing"
	"time"
)

const SqlitePath = "test.db"

type GormUserDataStoreTestSuite struct {
	suite.Suite
	db      *gorm.DB
	store   *GormUserDataStore
	user    *model.User
	userRow *sqlmock.Rows
}

func createUser(provider string) *model.User {
	return &model.User{
		ID:        uint(rand.Uint32()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     faker.Email(),
		Username:  faker.Username(),
		Provider:  provider,
		AvatarURL: faker.URL(),
		Files:     nil,
		Images:    nil,
		APIKey:    faker.UUIDDigit(),
	}
}

func TestGormUserDataStore(t *testing.T) {
	suite.Run(t, new(GormUserDataStoreTestSuite))
}

func (s *GormUserDataStoreTestSuite) SetupTest() {
	gormDB, err := gorm.Open(sqlite.Open(SqlitePath))
	require.NoError(s.T(), err)
	s.db = gormDB

	require.NoError(s.T(), gormDB.AutoMigrate(&model.User{}))

	s.store = &GormUserDataStore{db: gormDB}
	require.NoError(s.T(), err)

	s.user = createUser("github")
	require.NoError(s.T(), gormDB.Create(s.user).Error)
	require.NoError(s.T(), gormDB.Create(createUser("github")).Error)
	require.NoError(s.T(), gormDB.Create(createUser("google")).Error)
}

func (s *GormUserDataStoreTestSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), os.Remove(SqlitePath))
}

func (s *GormUserDataStoreTestSuite) TestFoundByID() {
	foundUser, err := s.store.FindById(s.user.ID)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(foundUser, s.user))
}

func (s *GormUserDataStoreTestSuite) TestNotFoundByID() {
	newUser := createUser("github")
	foundUser, err := s.store.FindById(newUser.ID)
	require.NoError(s.T(), err)
	require.Nil(s.T(), foundUser)
}

func (s *GormUserDataStoreTestSuite) TestFoundByProviderAndEmail() {
	foundUser, err := s.store.FindByProviderAndEmail(s.user.Provider, s.user.Email)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(foundUser, s.user))
}

func (s *GormUserDataStoreTestSuite) TestNotFoundByProviderAndEmail() {
	newUser := createUser("github")
	foundUser, err := s.store.FindByProviderAndEmail(newUser.Provider, newUser.Email)
	require.NoError(s.T(), err)
	require.Nil(s.T(), foundUser)
}

func (s *GormUserDataStoreTestSuite) TestFoundByAPIKey() {
	foundUser, err := s.store.FindByAPIKey(s.user.APIKey)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(foundUser, s.user))
}

func (s *GormUserDataStoreTestSuite) TestNotFoundByAPIKey() {
	newUser := createUser("github")
	foundUser, err := s.store.FindByAPIKey(newUser.APIKey)
	require.NoError(s.T(), err)
	require.Nil(s.T(), foundUser)
}

func (s *GormUserDataStoreTestSuite) TestUpdateAPIKey() {
	newApiKey := faker.UUIDDigit()
	require.NoError(s.T(), s.store.UpdateAPIKey(s.user.ID, newApiKey))

	queryUser := &model.User{}
	err := s.db.Where(&model.User{APIKey: newApiKey}).First(queryUser).Error
	require.NoError(s.T(), err)
	require.Equal(s.T(), queryUser.ID, s.user.ID)
}

func (s *GormUserDataStoreTestSuite) TestUpdateAPIKeyOnNotFoundUser() {
	newUser := createUser("github")
	newApiKey := faker.UUIDDigit()
	require.NoError(s.T(), s.store.UpdateAPIKey(newUser.ID, newApiKey))

	queryUser := &model.User{}
	err := s.db.Where(&model.User{APIKey: newApiKey}).First(queryUser).Error
	require.ErrorIs(s.T(), err, gorm.ErrRecordNotFound)
}
