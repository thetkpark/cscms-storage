package service

import (
	"github.com/matoous/go-nanoid/v2"
)

func GenerateFileToken() (string, error) {
	return gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 6)
}

func GenerateFileId() (string, error) {
	return gonanoid.New(30)
}

func GenerateImageToken() (string, error) {
	return gonanoid.Generate("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 20)
}
