package service

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.uber.org/zap"
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
	log  *zap.SugaredLogger
	path string
}

func NewDiskStorageManager(log *zap.SugaredLogger, path string) (*DiskStorageManager, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			log.Errorw("cannot create directory for file storage", "error", err)
			return nil, err
		}
	}
	if err != nil {
		log.Errorw("cannot check exist directory", "error", err)
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
		m.log.Errorw("cannot open file on disk", "error", err)
		return nil, err
	}
	return file, nil
}

func (m *DiskStorageManager) WriteToNewFile(fileName string, reader io.Reader) error {
	file, err := os.Create(m.getFilePath(fileName))
	if err != nil {
		m.log.Errorw("cannot create new file on disk", "error", err)
		return err
	}

	if _, err = io.Copy(file, reader); err != nil {
		m.log.Errorw("unable to write data to file", "error", err)
		return err
	}

	// Clean up the file
	if err := file.Sync(); err != nil {
		m.log.Errorw("unable sync the file to disk", "error", err)
		return err
	}
	if err := file.Close(); err != nil {
		m.log.Errorw("unable close the file", "error", err)
		return err
	}

	return nil
}

func (m *DiskStorageManager) Exist(fileName string) (bool, error) {
	if _, err := os.Stat(m.getFilePath(fileName)); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		m.log.Errorw("unable to check if file exist", "error", err)
		return false, err
	}
	return true, nil
}

func (m *DiskStorageManager) ListFiles() ([]string, error) {
	dir, err := os.Open(m.path)
	if err != nil {
		m.log.Errorw("unable to open storage directory", "error", err)
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		m.log.Errorw("unable to read file name in storage directory", "error", err)
		return nil, err
	}

	return files, nil
}

func (m *DiskStorageManager) DeleteFile(fileName string) error {
	if err := os.Remove(m.getFilePath(fileName)); err != nil {
		m.log.Errorw(fmt.Sprintf("unable to delete file %s", fileName), "error", err)
		return err
	}
	return nil
}

type ImageStorageManager interface {
	UploadImage(fileName string, file io.ReadSeekCloser) error
	DeleteImage(fileName string) error
}

type AzureImageStorageManager struct {
	log             *zap.SugaredLogger
	containerClient azblob.ContainerClient
}

func NewAzureImageStorageManager(l *zap.SugaredLogger, connectionString string, containerName string) (*AzureImageStorageManager, error) {
	serviceClient, err := azblob.NewServiceClientFromConnectionString(connectionString, nil)
	if err != nil {
		l.Errorw("Unable to create azblob service client", "error", err)
		return nil, err
	}
	containerClient := serviceClient.NewContainerClient(containerName)
	return &AzureImageStorageManager{
		log:             l,
		containerClient: containerClient,
	}, nil
}

func (a *AzureImageStorageManager) UploadImage(fileName string, file io.ReadSeekCloser) error {
	bbClient := a.containerClient.NewBlockBlobClient(fileName)
	_, err := bbClient.UploadStreamToBlockBlob(context.Background(), file, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		a.log.Errorw("Failed to upload file to az blob", "error", err)
	}
	return err
}

func (a *AzureImageStorageManager) DeleteImage(fileName string) error {
	bbClient := a.containerClient.NewBlockBlobClient(fileName)
	_, err := bbClient.Delete(context.Background(), nil)
	if err != nil {
		a.log.Errorw("Failed to delete file on az blob", "error", err)
	}
	return err
}
