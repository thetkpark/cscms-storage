package model

import (
	"gorm.io/gorm"
	"time"
)

type File struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Token     string    `gorm:"index" json:"token"`
	Nonce     string    `json:"nonce"`
	Filename  string    `json:"filename"`
	FileSize  uint64    `json:"file_size"`
	Visited   uint      `json:"visited"`
	UserID    uint      `gorm:"index"`
	FileType  string    `json:"file_type"`
	Encrypted bool      `json:"encrypted"`
	DeletedAt gorm.DeletedAt
}
