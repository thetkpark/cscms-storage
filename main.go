package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/handlers"
	"github.com/thetkpark/cscms-temp-storage/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thanhpk/randstr"
)

func main() {
	logger := hclog.Default()
	app := fiber.New(fiber.Config{
		BodyLimit: 150 << 20,
	})

	// Create data store
	db, err := gorm.Open(sqlite.Open("fileStore.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("unable to open sqlite db", err)
	}
	gormFileDataStore, err := data.NewGormFileDataStore(logger, db)
	if err != nil {
		log.Fatalln("unable to run gorm migration", err)
	}

	sioEncryptionManager := service.NewSIOEncryptionManager(logger, randstr.String(30))

	fileHandler := handlers.NewFileRoutesHandler(logger, sioEncryptionManager, gormFileDataStore)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET POST",
	}))

	app.Get("/api/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":   true,
			"timestamp": time.Now(),
		})
	})

	app.Post("/api/file", fileHandler.UploadFile)

	app.Get("/:token", fileHandler.GetFile)

	app.Static("/", "./client/build")

	err = app.Listen(":5000")
	if err != nil {
		log.Fatalln("unable to start server", err)
	}
}
