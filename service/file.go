package service

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"io"
	"os"
)

type StorageManager interface {
	CreateFile(fileName string) (io.Writer, error)
	OpenFile(fileName string) (io.Reader, error)
	WriteToNewFile(fileName string, reader io.Reader) error
	Exist(fileName string) (bool, error)
}

type DiskStorageManager struct {
	log  hclog.Logger
	path string
}

func NewDiskStorageManager(log hclog.Logger, path string) *DiskStorageManager {
	return &DiskStorageManager{path: path, log: log}
}

func (m *DiskStorageManager) getFilePath(fileName string) string {
	return fmt.Sprintf("%s/%s", m.path, fileName)
}

func (m *DiskStorageManager) CreateFile(fileName string) (io.Writer, error) {
	file, err := os.Create(m.getFilePath(fileName))
	if err != nil {
		m.log.Error("cannot create new file on disk", err)
		return nil, err
	}
	return file, nil
}

func (m *DiskStorageManager) OpenFile(fileName string) (io.Reader, error) {
	file, err := os.Open(m.getFilePath(fileName))
	if err != nil {
		m.log.Error("cannot open file on disk", err)
		return nil, err
	}
	return file, nil
}

func (m *DiskStorageManager) WriteToNewFile(fileName string, reader io.Reader) error {
	file, err := m.CreateFile(fileName)
	if err != nil {
		return err
	}

	if _, err = io.Copy(file, reader); err != nil {
		m.log.Error("unable to write data to file", err)
		return err
	}

	return nil
}

func (m *DiskStorageManager) Exist(fileName string) (bool, error) {
	if _, err := os.Stat(m.getFilePath(fileName)); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		m.log.Error("unable to check if file exist", err)
		return false, err
	}
	return true, nil
}
