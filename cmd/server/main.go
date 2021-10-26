package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/handlers"
	"github.com/thetkpark/cscms-temp-storage/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	logger := hclog.Default()

	masterKey, storagePath, port, maxStoreDuration := getEnv()
	dbHost, dbPort, dbUsername, dbPassword, dbName := getDBEnv()
	logger.Info("DB env", dbHost, dbPort, dbUsername, dbPassword, dbName)

	app := fiber.New(fiber.Config{
		BodyLimit: 150 << 20,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default to 500
			code := fiber.StatusInternalServerError
			message := err.Error()

			// Check if error is fiber.Error type
			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
				message = e.Message
			}

			c.Status(code)

			return c.JSON(fiber.Map{
				"code":    code,
				"message": message,
			})
		},
	})

	// Create data store
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)
	dsn := fmt.Sprintf("root:%s@tcp(common-mariadb-srv:3306)/cscms-temp-storage?charset=utf8mb4&parseTime=True&loc=Local", dbPassword)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("unable to open sqlite db", err)
	}
	gormFileDataStore, err := data.NewGormFileDataStore(logger, db, maxStoreDuration)
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
	fileHandler := handlers.NewFileRoutesHandler(logger, sioEncryptionManager, gormFileDataStore, diskStorageManager, maxStoreDuration)

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

	app.Static("/", "./client/build")

	app.Get("/:token", fileHandler.GetFile)

	err = app.Listen(port)
	if err != nil {
		log.Fatalf("unable to start server on %s: %v", port, err)
	}
}

func getEnv() (string, string, string, time.Duration) {
	// Required env
	key := os.Getenv("MASTER_KEY")
	storagePath := os.Getenv("STORAGE_PATH")
	if len(key) == 0 || len(storagePath) == 0 {
		log.Fatalln("MASTER_KEY and STORAGE_PATH env must be defined")
	}

	// Optional env
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = fmt.Sprintf(":%d", 5000)
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	maxStoreDuration := os.Getenv("STORE_DURATION") // in days
	duration := time.Hour * 24 * 30
	if len(maxStoreDuration) != 0 {
		date, err := strconv.Atoi(maxStoreDuration)
		if err != nil {
			log.Fatalln("STORE_DURATION is not a valid number")
		}
		duration = time.Hour * 24 * time.Duration(date)
	}

	return key, storagePath, port, duration
}

func getDBEnv() (string, string, string, string, string) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_DATABASE")

	return host, port, username, password, dbname
}
