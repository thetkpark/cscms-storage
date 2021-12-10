package model

import "time"

type Image struct {
	ID        uint      `gorm:"primaryKey,autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Filename  string    `json:"filename"`
	FileSize  uint64    `json:"file_size"`
	FilePath  string    `json:"file_path"`
	Token     string    `json:"token"`
	UserID    uint
}
