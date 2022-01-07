package handlers

import (
	"context"
	"io"
	"mime/multipart"
)

type Context interface {
	UserContext() context.Context
	SetUserContext(ctx context.Context)
	BaseURL() string
	Get(key string, defaultValue ...string) string
	Set(key, value string)
	FormFile(key string) (*multipart.FileHeader, error)
	Query(key string, defaultValue ...string) string
	Params(key string, defaultValue ...string) string
	Status(code int) Context
	Redirect(location string) error
	JSON(v interface{}) error
	SendStream(stream io.Reader, size ...int) error
	SendStatus(code int) error
	Next() error
	Error(code int, message string, error error) error
}

//SetCookie(name string, value string, expires time.Time, secure, HTTPOnly bool, samesite string)
//GetCookie(key, defaultValue string) string
//UserContext() context.Context
//SetUserContext(ctx context.Context)
