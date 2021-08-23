package service

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"io"
	"os"
)

type StorageManager interface {
	OpenFile(fileName string) (io.Reader, error)
	WriteToNewFile(fileName string, reader io.Reader) error
	Exist(fileName string) (bool, error)
	ListFiles() ([]string, error)
	DeleteFile(fileName string) error
}

type DiskStorageManager struct {
	log  hclog.Logger
	path string
}

func NewDiskStorageManager(log hclog.Logger, path string) (*DiskStorageManager, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			log.Error("cannot create directory for file storage", err)
			return nil, err
		}
	}
	if err != nil {
		log.Error("cannot check exist directory", err)
		return nil, err
	}

	return &DiskStorageManager{path: path, log: log}, nil
}

func (m *DiskStorageManager) getFilePath(fileName string) string {
	return fmt.Sprintf("%s/%s", m.path, fileName)
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
	file, err := os.Create(m.getFilePath(fileName))
	if err != nil {
		m.log.Error("cannot create new file on disk", err)
		return err
	}

	if _, err = io.Copy(file, reader); err != nil {
		m.log.Error("unable to write data to file", err)
		return err
	}

	// Clean up the file
	if err := file.Close(); err != nil {
		m.log.Error("unable close the file", err)
		return err
	}
	if err := file.Sync(); err != nil {
		m.log.Error("unable sync the file to disk", err)
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

func (m *DiskStorageManager) ListFiles() ([]string, error) {
	dir, err := os.Open(m.path)
	if err != nil {
		m.log.Error("unable to open storage directory", err)
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		m.log.Error("unable to read file name in storage directory", err)
		return nil, err
	}

	return files, nil
}

func (m *DiskStorageManager) DeleteFile(fileName string) error {
	if err := os.Remove(m.getFilePath(fileName)); err != nil {
		m.log.Error(fmt.Sprintf("unable to delete file %s", fileName), err)
		return err
	}
	return nil
}
