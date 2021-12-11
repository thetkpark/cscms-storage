package data

import (
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"gorm.io/gorm"
)

type ImageDataStore interface {
	Create(image *model.Image) error
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

func (g *GormImageDataStore) Create(image *model.Image) error {
	tx := g.db.Create(image)
	return tx.Error
}
