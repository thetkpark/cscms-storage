package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/thanhpk/randstr"
	"github.com/thetkpark/cscms-temp-storage/service"
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

var EncryptionKey = randstr.Bytes(32)

func (h *FileRoutesHandler) UploadFile(c *fiber.Ctx) error {
	return nil
}

func (h *FileRoutesHandler) GetFile(c *fiber.Ctx) error {
	return nil
}
