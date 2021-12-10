package data

import (
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"gorm.io/gorm"
)

type UserDataStore interface {
}

type GormUserDataStore struct {
	log hclog.Logger
	db  *gorm.DB
}

func NewGormUserDataStore(l hclog.Logger, db *gorm.DB) (*GormUserDataStore, error) {
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}
	return &GormUserDataStore{
		log: l,
		db:  db,
	}, nil
}
