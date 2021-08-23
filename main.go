package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/handlers"
	"github.com/thetkpark/cscms-temp-storage/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	logger := hclog.Default()

	masterKey, storagePath, sqlitePath := getEnv()

	app := fiber.New(fiber.Config{
		BodyLimit: 150 << 20,
	})

	// Create data store
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		log.Fatalln("unable to open sqlite db", err)
	}
	gormFileDataStore, err := data.NewGormFileDataStore(logger, db)
	if err != nil {
		log.Fatalln("unable to run gorm migration", err)
	}

	// Create service managers for handler
	sioEncryptionManager := service.NewSIOEncryptionManager(logger, masterKey)
	diskStorageManager, err := service.NewDiskStorageManager(logger, storagePath)
	if err != nil {
		log.Fatalln("unable to create disk storage manager")
	}

	// Create handlers
	fileHandler := handlers.NewFileRoutesHandler(logger, sioEncryptionManager, gormFileDataStore, diskStorageManager)

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

func getEnv() (string, string, string) {
	key := os.Getenv("MASTER_KEY")
	storagePath := os.Getenv("STORAGE_PATH")
	dbPath := os.Getenv("SQLITE_PATH")
	if len(key) == 0 || len(storagePath) == 0 || len(dbPath) == 0 {
		log.Fatalln("MASTER_KEY, SQLITE_PATH STORAGE_PATH env must be defined")
	}
	return key, storagePath, dbPath
}
