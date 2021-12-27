package storage

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"testing"
)

const StoragePath = "files"

func createDiskStorageManager() (*DiskStorageManager, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	sugarLogger := logger.Sugar()
	return NewDiskStorageManager(sugarLogger, StoragePath)
}

func cleanup() error {
	return os.RemoveAll(StoragePath)
}

func createTestFile(fileName string, fileContent string) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", StoragePath, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	fileContentReader := strings.NewReader(fileContent)
	_, err = io.Copy(file, fileContentReader)
	return err
}

func TestWriteToNewFile(t *testing.T) {
	diskStorageManager, err := createDiskStorageManager()
	require.NoError(t, err)
	defer cleanup()

	fileName := "test-write"
	fileContent := "hello world. This is some string to be tested"
	fileInputReader := strings.NewReader(fileContent)
	err = diskStorageManager.WriteToNewFile(fileName, fileInputReader)
	require.NoError(t, err)

	// Check file on disk
	diskFile, err := os.Open(fmt.Sprintf("%s/%s", StoragePath, fileName))
	require.NoError(t, err)
	diskFileString := new(strings.Builder)
	_, err = io.Copy(diskFileString, diskFile)
	require.NoError(t, err)
	require.Equal(t, fileContent, diskFileString.String())
}

func TestOpenFile(t *testing.T) {
	diskStorageManager, err := createDiskStorageManager()
	require.NoError(t, err)
	defer cleanup()

	fileName := "test-open-file"
	fileContent := "When I was a young boy, my father took me into the city to see a marching band"
	err = createTestFile(fileName, fileContent)
	require.NoError(t, err)

	fileReader, err := diskStorageManager.OpenFile(fileName)
	require.NoError(t, err)

	// Checking
	fileString := new(strings.Builder)
	_, err = io.Copy(fileString, fileReader)
	require.NoError(t, err)
	require.Equal(t, fileString.String(), fileContent)
}

func TestOpenDeletedFile(t *testing.T) {
	diskStorageManager, err := createDiskStorageManager()
	require.NoError(t, err)
	defer cleanup()

	_, err = diskStorageManager.OpenFile("doesNotExist")
	require.Error(t, err)
}

func TestFileExist(t *testing.T) {
	diskStorageManager, err := createDiskStorageManager()
	require.NoError(t, err)
	defer cleanup()

	fileName := "test-file-exist"
	fileContent := "When I was a young boy, my father took me into the city to see a marching band"
	err = createTestFile(fileName, fileContent)
	require.NoError(t, err)

	isExist, err := diskStorageManager.Exist(fileName)
	require.NoError(t, err)
	require.True(t, isExist)
}

func TestFileNotExist(t *testing.T) {
	diskStorageManager, err := createDiskStorageManager()
	require.NoError(t, err)
	defer cleanup()

	isExist, err := diskStorageManager.Exist("doesNotExist")
	require.NoError(t, err)
	require.False(t, isExist)
}

func TestListFiles(t *testing.T) {
	diskStorageManager, err := createDiskStorageManager()
	require.NoError(t, err)
	defer cleanup()

	fileNameLists := []string{"test-list-1", "test-list-2", "test-list-3", "test-list-4"}
	fileContent := "When I was a young boy, my father took me into the city to see a marching band"
	for _, fileName := range fileNameLists {
		err := createTestFile(fileName, fileContent)
		require.NoError(t, err)
	}

	fileLists, err := diskStorageManager.ListFiles()
	require.NoError(t, err)
	require.EqualValues(t, fileLists, fileNameLists)
}

func TestDeleteFile(t *testing.T) {
	diskStorageManager, err := createDiskStorageManager()
	require.NoError(t, err)
	defer cleanup()

	fileName := "test-file-exist"
	fileContent := "When I was a young boy, my father took me into the city to see a marching band"
	err = createTestFile(fileName, fileContent)
	require.NoError(t, err)

	err = diskStorageManager.DeleteFile(fileName)
	require.NoError(t, err)

	_, err = os.Open(fmt.Sprintf("%s/%s", StoragePath, fileName))
	require.Error(t, err)
}
