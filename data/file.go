package data

import (
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"gorm.io/gorm"
	"time"
)

type FileDataStore interface {
	Create(file *model.File) error
	FindByToken(token string) (*model.File, error)
	IncreaseVisited(id string) error
	FindByUserID(userId uint) (*[]model.File, error)
	FindByUserIDAndFileID(userId uint, fileId string) (*model.File, error)
	DeleteByID(fileId string) error
}

type GormFileDataStore struct {
	log              hclog.Logger
	db               *gorm.DB
	maxStoreDuration time.Duration
}

func NewGormFileDataStore(l hclog.Logger, db *gorm.DB, duration time.Duration) (*GormFileDataStore, error) {
	if err := db.AutoMigrate(&model.File{}); err != nil {
		return nil, err
	}

	return &GormFileDataStore{
		log:              l,
		db:               db,
		maxStoreDuration: duration,
	}, nil
}

func (store *GormFileDataStore) Create(file *model.File) error {
	tx := store.db.Create(file)
	return tx.Error
}

func (store *GormFileDataStore) FindByToken(token string) (*model.File, error) {
	var files []*model.File
	tx := store.db.Where(&model.File{Token: token}).Find(&files)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			store.log.Error("error querying file by token")
		}
		return nil, tx.Error
	}

	var file *model.File
	for _, v := range files {
		if v.ExpiredAt.UTC().After(time.Now().UTC()) {
			file = v
			break
		}
	}

	if file == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return file, nil
}

func (store *GormFileDataStore) FindByUserIDAndFileID(userId uint, fileId string) (*model.File, error) {
	var file model.File
	tx := store.db.Where(&model.File{ID: fileId, UserID: userId}).Limit(1).Find(&file)
	return &file, tx.Error
}

func (store *GormFileDataStore) FindByUserID(userId uint) (*[]model.File, error) {
	var files []model.File
	tx := store.db.Where(&model.File{UserID: userId}).Find(&files)
	return &files, tx.Error
}

func (store *GormFileDataStore) DeleteByID(fileId string) error {
	tx := store.db.Delete(&model.File{ID: fileId})
	return tx.Error
}

func (store *GormFileDataStore) IncreaseVisited(id string) error {
	tx := store.db.Table("files").Where(&model.File{ID: id}).UpdateColumns(map[string]interface{}{
		"visited":    gorm.Expr("visited + ?", 1),
		"updated_at": time.Now().UTC(),
	})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
