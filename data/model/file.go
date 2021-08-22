package model

import "time"

type File struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time
	Token     string `gorm:"unique uniqueIndex" json:"token"`
	Nonce     string `json:"nonce"`
	Filename  string `json:"filename"`
	FileSize  uint64 `json:"file_size"`
	Visited   uint   `json:"visited"`
}