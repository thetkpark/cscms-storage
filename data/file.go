package data

import (
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"gorm.io/gorm"
	"time"
)

type FileDataStore interface {
	Create(id, token, nonce, filename string, filesize uint64, storeDuration time.Duration) (*model.File, error)
	FindByToken(token string) (*model.File, error)
	IncreaseVisited(id string) error
}

type GormFileDataStore struct {
	log hclog.Logger
	db  *gorm.DB
}

func NewGormFileDataStore(l hclog.Logger, db *gorm.DB) (*GormFileDataStore, error) {
	if err := db.AutoMigrate(&model.File{}); err != nil {
		return nil, err
	}

	return &GormFileDataStore{
		log: l,
		db:  db,
	}, nil
}

func (store *GormFileDataStore) Create(id, token, nonce, filename string, filesize uint64, storeDuration time.Duration) (*model.File, error) {
	file := &model.File{
		ID:        id,
		Token:     token,
		Nonce:     nonce,
		Filename:  filename,
		FileSize:  filesize,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiredAt: time.Now().UTC().Add(storeDuration),
		Visited:   0,
	}

	tx := store.db.Create(file)
	if tx.Error != nil {
		store.log.Error("unable to create new file data in db", tx.Error)
		return nil, tx.Error
	}
	return file, nil
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
