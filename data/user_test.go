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

func createUser() *model.User {
	return &model.User{
		ID:        uint(rand.Uint32()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     faker.Email(),
		Username:  faker.Username(),
		Provider:  "GitHub",
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

	require.NoError(s.T(), gormDB.AutoMigrate(&model.User{}))

	s.store = &GormUserDataStore{db: gormDB}
	require.NoError(s.T(), err)

	s.user = createUser()
	require.NoError(s.T(), gormDB.Create(s.user).Error)
	require.NoError(s.T(), gormDB.Create(createUser()).Error)
	require.NoError(s.T(), gormDB.Create(createUser()).Error)
}

func (s *GormUserDataStoreTestSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), os.Remove(SqlitePath))
}

func (s *GormUserDataStoreTestSuite) TestFindByID() {
	foundUser, err := s.store.FindById(s.user.ID)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(foundUser, s.user))
}
