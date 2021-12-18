package model

import (
	"gorm.io/gorm"
	"time"
)

type Image struct {
	ID               uint      `gorm:"primaryKey,autoIncrement" json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	OriginalFilename string    `json:"original_filename"`
	FileSize         uint64    `json:"file_size"`
	FilePath         string    `json:"file_path"`
	UserID           uint      `json:"user_id" gorm:"index"`
	DeletedAt        gorm.DeletedAt
}
