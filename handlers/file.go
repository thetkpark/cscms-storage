package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"github.com/thetkpark/cscms-temp-storage/service"
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
		"total_time":   time.Since(tStart).String(),
	})
}

func (h *FileRoutesHandler) GetFile(c *fiber.Ctx) error {
	tStart := time.Now()
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

	// Check if file still exist on storage
	if exist, err := h.storageManager.Exist(file.ID); !exist {
		if err == nil {
			// File is not exist anymore
			return fiber.NewError(fiber.StatusNotFound, "file not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "unable to check if file exist", err.Error())
	}

	// Get encrypted file from storage manager
	encryptedFile, err := h.storageManager.OpenFile(file.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to open encrypted file", err.Error())
	}

	// Decrypt file
	tEnc := time.Now()
	err = h.encryptionManager.Decrypt(encryptedFile, file.Nonce, c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to decrypt", err.Error())
	}
	decryptDuration := time.Since(tEnc)

	// Increase visited count
	err = h.fileDataStore.IncreaseVisited(file.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to increase count", err.Error())
	}

	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.Filename))
	c.Set("Content-Type", "application/octet-stream")

	h.log.Info(file.ID)
	h.log.Info("decrypt duration", decryptDuration.String())
	h.log.Info("total duration", time.Since(tStart).String())
	return nil
}
