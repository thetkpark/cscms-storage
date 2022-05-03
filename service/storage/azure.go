package storage

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.uber.org/zap"
	"io"
)

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

func (a *AzureImageStorageManager) UploadImage(fileName, mimeType string, file io.ReadSeekCloser) error {
	bbClient := a.containerClient.NewBlockBlobClient(fileName)
	_, err := bbClient.Upload(context.Background(), file, &azblob.UploadBlockBlobOptions{
		HTTPHeaders: &azblob.BlobHTTPHeaders{ BlobContentType: &mimeType },
	})
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
