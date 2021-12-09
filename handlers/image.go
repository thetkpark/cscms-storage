package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/service"
	"regexp"
)

type ImageRouteHandler struct {
	log               hclog.Logger
	imageStoreManager service.ImageStorageManager
}

func NewImageRouteHandler(log hclog.Logger, store service.ImageStorageManager) *ImageRouteHandler {
	return &ImageRouteHandler{
		log:               log,
		imageStoreManager: store,
	}
}

func (h *ImageRouteHandler) UploadImage(c *fiber.Ctx) error {
	// Get image file from Form
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return NewHTTPError(h.log, fiber.StatusBadRequest, "Unable to get file from multipart/form-data", err)
	}

	// Check image format and get extension
	fileExtension, err := h.validateFileFormat(fileHeader.Header.Get("Content-Type"), fileHeader.Filename)
	if err != nil {
		return NewHTTPError(h.log, fiber.StatusBadRequest, "Invalid file extension", err)
	}
	imageToken, err := service.GenerateImageToken()
	if err != nil {
		return NewHTTPError(h.log, fiber.StatusBadRequest, "Unable to generate image token", err)
	}
	imagePath := fmt.Sprintf("%s.%s", imageToken, fileExtension)

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return NewHTTPError(h.log, fiber.StatusInternalServerError, "Unable to open file from the fileHeader", err)
	}

	// Upload the image to storage
	if err := h.imageStoreManager.UploadImage(imagePath, file); err != nil {
		return NewHTTPError(h.log, fiber.StatusInternalServerError, "Unable to upload image", err)
	}

	// Return the
	return c.JSON(fiber.Map{
		"image_path": imagePath,
	})
}

func (h *ImageRouteHandler) validateFileFormat(mimeType string, fileName string) (string, error) {
	switch mimeType {
	case "image/png", "image/jpeg", "image/gif", "image/x-icon", "image/heic", "image/webp", "image/tiff", "image/svg+xml", "image/bmp", "image/apng", "image/avif":
		return h.getFileExtension(fileName), nil
	default:
		return "", fmt.Errorf("%s is not supported", mimeType)
	}

}

func (h *ImageRouteHandler) getFileExtension(fileName string) string {
	regex := regexp.MustCompile(".+\\.(.+)")
	group := regex.FindStringSubmatch(fileName)
	return group[1]
}
