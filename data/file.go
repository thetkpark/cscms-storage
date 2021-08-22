package data

import (
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"gorm.io/gorm"
)

type FileDataStore interface {
	Create(id, token, nonce, filename string, filesize uint64) (*model.File, error)
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

func (store *GormFileDataStore) Create(id, token, nonce, filename string, filesize uint64) (*model.File, error) {
	file := &model.File{
		ID:       id,
		Token:    token,
		Nonce:    nonce,
		Filename: filename,
		FileSize: filesize,
		Visited:  0,
	}

	tx := store.db.Create(file)
	if tx.Error != nil {
		store.log.Error("unable to create new file data in db", tx.Error)
		return nil, tx.Error
	}
	return file, nil
}

func (store *GormFileDataStore) FindByToken(token string) (*model.File, error) {
	file := &model.File{}
	tx := store.db.Where(&model.File{Token: token}).First(file)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			store.log.Error("error querying file by token")
		}
		return nil, tx.Error
	}
	return file, nil
}

func (store *GormFileDataStore) IncreaseVisited(id string) error {
	tx := store.db.Where(&model.File{ID: id}).Update("visited", gorm.Expr("visited + ?", 1))
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
