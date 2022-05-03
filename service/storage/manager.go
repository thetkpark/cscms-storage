package storage

import "io"

type FileManager interface {
	OpenFile(fileName string) (io.Reader, error)
	WriteToNewFile(fileName string, reader io.Reader) error
	Exist(fileName string) (bool, error)
	ListFiles() ([]string, error)
	DeleteFile(fileName string) error
}

type ImageManager interface {
	UploadImage(fileName string, mimeType string,file io.ReadSeekCloser) error
	DeleteImage(fileName string) error
}
