package service

import (
	"github.com/matoous/go-nanoid/v2"
)

func GenerateFileToken() (string, error) {
	token, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 6)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateFileId() (string, error) {
	key, err := gonanoid.New(30)
	if err != nil {
		return "", err
	}
	return key, nil
}
