package handlers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/service"
	"gorm.io/gorm"
	"time"
)

type FileRoutesHandler struct {
	log               hclog.Logger
	encryptionManager service.EncryptionManager
	fileDataStore     data.FileDataStore
	storageManager    service.StorageManager
}

func NewFileRoutesHandler(log hclog.Logger, enc service.EncryptionManager, data data.FileDataStore, store service.StorageManager) *FileRoutesHandler {
	return &FileRoutesHandler{
		log:               log,
		encryptionManager: enc,
		fileDataStore:     data,
		storageManager:    store,
	}
}

func (h *FileRoutesHandler) UploadFile(c *fiber.Ctx) error {
	tStart := time.Now()
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to get file from form-data", err.Error())
	}

	// Create new fileId
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
	tEncrypt := time.Now()
	encrypted, nonce, err := h.encryptionManager.Encrypt(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable encrypt the file", err.Error())
	}
	encryptFileDuration := time.Now().Sub(tEncrypt)

	// Write encrypted file to disk
	if err := h.storageManager.WriteToNewFile(fileId, encrypted); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to write encrypted data to file", err.Error())
	}

	// Check slug
	fileToken := c.Query("slug")
	if len(fileToken) == 0 {
		//No slug, Generate file token
		fileToken, err = service.GenerateFileToken()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "unable to generate file token", err.Error())
		}
	} else {
		// Check if slug is available
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
		"total_time":   time.Since(tStart).String(),
	})
}

func (h *FileRoutesHandler) GetFile(c *fiber.Ctx) error {
	token := c.Params("token")

	// Find file by token
	file, err := h.fileDataStore.FindByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Redirect(c.BaseURL())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "unable to get file query", err.Error())
	}

	// Check if file still exist on storage
	if exist, err := h.storageManager.Exist(file.ID); !exist {
		if err == nil {
			// File is not exist anymore
			return c.Redirect(c.BaseURL())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "unable to check if file exist", err.Error())
	}

	// Get encrypted file from storage manager
	encryptedFile, err := h.storageManager.OpenFile(file.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to open encrypted file", err.Error())
	}

	// Decrypt file
	err = h.encryptionManager.Decrypt(encryptedFile, file.Nonce, c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to decrypt", err.Error())
	}

	// Increase visited count
	err = h.fileDataStore.IncreaseVisited(file.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to increase count", err.Error())
	}

	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.Filename))
	c.Set("Content-Type", "application/octet-stream")

	return nil
}
