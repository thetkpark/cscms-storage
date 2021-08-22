package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/service"
	"io"
	"os"
	"runtime"
	"time"
)

type FileRoutesHandler struct {
	log               hclog.Logger
	encryptionManager service.EncryptionManager
}

func NewFileRoutesHandler(log hclog.Logger, encryptManager service.EncryptionManager) *FileRoutesHandler {
	return &FileRoutesHandler{
		log:               log,
		encryptionManager: encryptManager,
	}
}

func (h *FileRoutesHandler) UploadFile(c *fiber.Ctx) error {
	fmt.Println("Start")
	PrintMemUsage()
	t1 := time.Now()
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to get file from form-data", err.Error())
	}

	// Create new random uuid
	fileId, err := uuid.NewRandom()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to create new filename", err.Error())
	}

	// Open multipart form header
	fmt.Println("Before open file")
	PrintMemUsage()
	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to open file", err.Error())
	}
	fmt.Println("After open file")
	PrintMemUsage()

	// Encrypt the file
	fmt.Println("Before encrypt file")
	PrintMemUsage()
	ts := time.Now()
	encrypted, key, err := h.encryptionManager.Encrypt(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable encrypt the file", err.Error())
	}
	encryptFileDuration := time.Now().Sub(ts)
	fmt.Println("After encrypt file")
	PrintMemUsage()

	// Create file
	encryptedFilePath := fmt.Sprintf("%s/%s.enc", "tmp", fileId)
	encryptedFile, err := os.Create(encryptedFilePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to create new file on disk", err.Error())
	}
	defer encryptedFile.Close()

	// Write encrypted file to disk
	fmt.Println("Before write file to disk")
	PrintMemUsage()
	if _, err := io.Copy(encryptedFile, encrypted); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to write bytes to file", err.Error())
	}
	fmt.Println("After write file to disk")
	PrintMemUsage()

	return c.JSON(fiber.Map{
		"id":               fileId,
		"key":              key,
		"encrypt_duration": encryptFileDuration.String(),
		"total_time":       time.Since(t1).String(),
	})
}

func (h *FileRoutesHandler) GetFile(c *fiber.Ctx) error {
	return nil
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
