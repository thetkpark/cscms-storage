package data

import (
	"errors"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"gorm.io/gorm"
)

type UserDataStore interface {
	FindByProviderAndEmail(provider string, email string) (*model.User, error)
	FindById(userId uint) (*model.User, error)
	Create(email string, username string, provider string, avatarUrl string) (*model.User, error)
}

type GormUserDataStore struct {
	db *gorm.DB
}

func NewGormUserDataStore(db *gorm.DB) (*GormUserDataStore, error) {
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}
	return &GormUserDataStore{
		db: db,
	}, nil
}

func (d *GormUserDataStore) FindById(userId uint) (*model.User, error) {
	var user model.User
	tx := d.db.Where(&model.User{ID: userId}).First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &user, nil
}

func (d *GormUserDataStore) FindByProviderAndEmail(provider string, email string) (*model.User, error) {
	var user model.User
	tx := d.db.Where(&model.User{Email: email, Provider: provider}).First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &user, nil
}

func (d *GormUserDataStore) Create(email string, username string, provider string, avatarUrl string) (*model.User, error) {
	user := &model.User{
		Email:     email,
		Username:  username,
		Provider:  provider,
		AvatarURL: avatarUrl,
	}
	tx := d.db.Create(user)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
