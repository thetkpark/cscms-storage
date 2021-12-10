package data

import (
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"gorm.io/gorm"
)

type ImageDataStore interface {
}

type GormImageDataStore struct {
	log hclog.Logger
	db  *gorm.DB
}

func NewGormImageDataStore(l hclog.Logger, db *gorm.DB) (*GormImageDataStore, error) {
	if err := db.AutoMigrate(&model.Image{}); err != nil {
		return nil, err
	}
	return &GormImageDataStore{
		log: l,
		db:  db,
	}, nil
}
