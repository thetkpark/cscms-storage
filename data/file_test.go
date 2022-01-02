package data

import (
	"github.com/bxcodec/faker/v3"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"testing"
	"time"
)

type GormFileDataStoreTestSuite struct {
	suite.Suite
	db      *gorm.DB
	store   *GormFileDataStore
	user    *model.User
	file    *model.File
	ownFile *[]model.File
}

func createTestFile(userID uint, expired bool) *model.File {
	file := &model.File{
		ID:        faker.Password(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Hour),
		Token:     faker.Password(),
		Nonce:     faker.UUIDDigit(),
		Filename:  faker.Username(),
		FileSize:  uint64(rand.Uint32()),
		Visited:   uint(rand.Uint32()),
		UserID:    0,
		FileType:  faker.Currency(),
		Encrypted: rand.Int() > rand.Int(),
	}
	if userID != 0 {
		file.UserID = userID
	}
	if expired {
		file.ExpiredAt = time.Now().Add(-1 * time.Hour)
	}
	return file
}

func TestGormFileDataStore(t *testing.T) {
	suite.Run(t, new(GormFileDataStoreTestSuite))
}

func (s *GormFileDataStoreTestSuite) SetupTest() {
	gormDB, err := createGormDB()
	require.NoError(s.T(), err)
	s.db = gormDB

	require.NoError(s.T(), gormDB.AutoMigrate(&model.File{}, &model.User{}))

	s.store = &GormFileDataStore{db: gormDB}
	require.NoError(s.T(), err)

	s.user = createTestUser("github")
	s.file = createTestFile(0, false)
	s.ownFile = &[]model.File{
		*createTestFile(s.user.ID, false),
		*createTestFile(s.user.ID, false),
		*createTestFile(s.user.ID, true),
	}
	require.NoError(s.T(), s.db.Create(s.user).Error)
	require.NoError(s.T(), s.db.Create(s.file).Error)
	require.NoError(s.T(), s.db.Create(&s.ownFile).Error)
}

func (s *GormFileDataStoreTestSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), os.Remove(SqlitePath))
}

func (s *GormFileDataStoreTestSuite) TestCreate() {
	newFile := createTestFile(0, false)
	require.NoError(s.T(), s.store.Create(newFile))

	var queryFile model.File
	require.NoError(s.T(), s.db.Where(newFile).First(&queryFile).Error)
	require.Nil(s.T(), deep.Equal(&queryFile, newFile))
}

func (s *GormFileDataStoreTestSuite) TestFindByID() {
	file, err := s.store.FindByID(s.file.ID)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(file, s.file))
}

func (s *GormFileDataStoreTestSuite) TestFindByIDNotFound() {
	newFile := createTestFile(0, false)
	file, err := s.store.FindByID(newFile.ID)
	require.NoError(s.T(), err)
	require.Nil(s.T(), file)
}

func (s *GormFileDataStoreTestSuite) TestFindByToken() {
	file, err := s.store.FindByToken(s.file.Token)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(file, s.file))
}

func (s *GormFileDataStoreTestSuite) TestFindByTokenNotFound() {
	newFile := createTestFile(0, false)
	file, err := s.store.FindByToken(newFile.Token)
	require.NoError(s.T(), err)
	require.Nil(s.T(), file)
}

func (s *GormFileDataStoreTestSuite) TestFindByTokenExpired() {
	newFile := createTestFile(0, true)
	s.db.Create(newFile)

	file, err := s.store.FindByToken(newFile.Token)
	require.NoError(s.T(), err)
	require.Nil(s.T(), file)
}

func (s *GormFileDataStoreTestSuite) TestFindByUserID() {
	files, err := s.store.FindByUserID(s.user.ID)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(files, s.ownFile))
}

func (s *GormFileDataStoreTestSuite) TestFindByUserIDEmpty() {
	newUser := createTestUser("google")
	files, err := s.store.FindByUserID(newUser.ID)
	require.NoError(s.T(), err)
	require.Len(s.T(), *files, 0)
}

func (s *GormFileDataStoreTestSuite) TestIncreaseVisited() {
	require.NoError(s.T(), s.store.IncreaseVisited(s.file.ID))

	var queryFile model.File
	require.NoError(s.T(), s.db.Where("id", s.file.ID).First(&queryFile).Error)
	require.Equal(s.T(), queryFile.Visited, s.file.Visited+1)
}

func (s *GormFileDataStoreTestSuite) TestDeleteByID() {
	require.NoError(s.T(), s.store.DeleteByID(s.file.ID))
	var queryFile model.File
	require.ErrorIs(s.T(), s.db.Where("id", s.file.ID).First(&queryFile).Error, gorm.ErrRecordNotFound)
}

func (s *GormFileDataStoreTestSuite) TestUpdateToken() {
	s.file.Token = "newToken"
	require.NoError(s.T(), s.store.UpdateToken(s.file.ID, s.file.Token))
	var queryFile model.File
	require.NoError(s.T(), s.db.Where("token", s.file.Token).First(&queryFile).Error)
	require.Nil(s.T(), deep.Equal(&queryFile, s.file))
}
