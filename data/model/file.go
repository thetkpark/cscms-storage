package model

import "time"

type File struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `gorm:"unique uniqueIndex" json:"token"`
	Filename  string    `json:"filename"`
	FileSize  uint      `json:"file_size"`
	Visited   uint      `json:"visited"`
}
