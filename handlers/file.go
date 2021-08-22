package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"github.com/thetkpark/cscms-temp-storage/service"
	"io"
	"os"
	"time"
)

type FileRoutesHandler struct {
	log               hclog.Logger
	encryptionManager service.EncryptionManager
	fileDataStore     data.FileDataStore
}

func NewFileRoutesHandler(log hclog.Logger, encryptManager service.EncryptionManager, dataStore data.FileDataStore) *FileRoutesHandler {
	return &FileRoutesHandler{
		log:               log,
		encryptionManager: encryptManager,
		fileDataStore:     dataStore,
	}
}

func (h *FileRoutesHandler) UploadFile(c *fiber.Ctx) error {
	t1 := time.Now()
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to get file from form-data", err.Error())
	}

	// Create new random uuid
	fileId, err := service.GenerateFileId()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to create file id", err.Error())
	}

	// Open multipart form header
	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to open file", err.Error())
	}

	// Encrypt the file
	ts := time.Now()
	encrypted, nonce, err := h.encryptionManager.Encrypt(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable encrypt the file", err.Error())
	}
	encryptFileDuration := time.Now().Sub(ts)

	// Create file
	encryptedFilePath := fmt.Sprintf("%s/%s", "tmp", fileId)
	encryptedFile, err := os.Create(encryptedFilePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to create new file on disk", err.Error())
	}
	defer encryptedFile.Close()

	// Write encrypted file to disk
	if _, err := io.Copy(encryptedFile, encrypted); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to write bytes to file", err.Error())
	}

	// Generate file token
	fileToken, err := service.GenerateFileToken()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to generate file token", err.Error())
	}

	// Create new fileInfo record in db
	fileInfo, err := h.fileDataStore.Create(fileId, fileToken, nonce, fileHeader.Filename, uint64(fileHeader.Size))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to save file info to db", err.Error())
	}

	return c.JSON(fiber.Map{
		"id":           fileInfo.ID,
		"token":        fileInfo.Token,
		"file_size":    fileInfo.FileSize,
		"file_name":    fileInfo.Filename,
		"created_at":   fileInfo.CreatedAt,
		"encrypt_time": encryptFileDuration.String(),
		"total_time":   time.Since(t1).String(),
	})
}

func (h *FileRoutesHandler) GetFile(c *fiber.Ctx) error {
	t1 := time.Now()
	token := c.Params("token")

	// Find file by token
	files, err := h.fileDataStore.FindByToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "unable to find a file info from db", err.Error())
	}

	// Looping to find the unexpired file
	var file *model.File
	for _, v := range files {
		if v.CreatedAt.UTC().Add(time.Hour * 720).After(time.Now().UTC()) {
			file = v
			break
		}
	}
	if file == nil {
		return fiber.NewError(fiber.StatusNotFound, "file not found")
	}

	encryptedFilePath := fmt.Sprintf("%s/%s", "tmp", file.ID)
	if _, err := os.Stat(encryptedFilePath); os.IsNotExist(err) {
		return fiber.NewError(fiber.StatusNotFound, "file not found")
	}

	// open encrypted file
	encryptedFile, err := os.Open(encryptedFilePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to open encrypted file", err.Error())
	}

	// Decrypt file
	ts := time.Now()
	err = h.encryptionManager.Decrypt(encryptedFile, file.Nonce, c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to decrypt", err.Error())
	}
	decryptDuration := time.Since(ts)

	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.Filename))
	c.Set("Content-Type", "application/octet-stream")

	h.log.Info(file.ID)
	h.log.Info("decrypt duration", decryptDuration.String())
	h.log.Info("total duration", time.Since(t1).String())
	return nil
}
