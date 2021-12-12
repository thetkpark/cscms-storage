package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/markbates/goth"
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
	//"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
)

func main() {
	logger := hclog.Default()

	masterKey, storagePath, port, maxStoreDuration, azStorageConnString, azStorageConName := getEnv()
	dbHost, dbPort, dbUsername, dbPassword, dbName := getDBEnv()
	entrypoint, ghClientId, ghSecretKey, ggClientId, ggSecretKey := getOauthEnv()

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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("unable to open sqlite db", err)
	}
	gormFileDataStore, err := data.NewGormFileDataStore(logger, db, maxStoreDuration)
	if err != nil {
		log.Fatalln("unable to run gorm migration on file table", err)
	}
	gormImageDataStore, err := data.NewGormImageDataStore(logger, db)
	if err != nil {
		log.Fatalln("unable to run gorm migration on image table", err)
	}
	gormUserDataStore, err := data.NewGormUserDataStore(logger, db)
	if err != nil {
		log.Fatalln("unable to run gorm migration on user table", err)
	}

	// Create service managers for handler
	sioEncryptionManager := service.NewSIOEncryptionManager(logger, masterKey)
	diskStorageManager, err := service.NewDiskStorageManager(logger, storagePath)
	if err != nil {
		log.Fatalln("unable to create disk storage manager")
	}
	imageStorageManager, err := service.NewAzureImageStorageManager(logger, azStorageConnString, azStorageConName)
	if err != nil {
		log.Fatalln("unable to azure image storage manager")
	}
	jwtManager := service.NewJwtManager(os.Getenv("JWT_SECRET"))

	// Create handlers
	fileHandler := handlers.NewFileRoutesHandler(logger, sioEncryptionManager, gormFileDataStore, diskStorageManager, maxStoreDuration)
	imageHandler := handlers.NewImageRouteHandler(logger, gormImageDataStore, imageStorageManager)
	authHandler := handlers.NewAuthRouteHandler(logger, gormUserDataStore, jwtManager, entrypoint)

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
	app.Post("/api/file", authHandler.ParseUserFromCookie, fileHandler.UploadFile)
	//app.Get("/api/file/lists", authHandler.ParseUserFromCookie, authHandler.AuthenticatedOnly)

	app.Post("/api/image", imageHandler.UploadImage)
	//app.Get("/api/image/lists", authHandler.ParseUserFromCookie, authHandler.AuthenticatedOnly)

	app.Static("/", "./client/build")
	app.Static("/404", "./client/build")

	// Try auth
	goth.UseProviders(
		github.New(ghClientId, ghSecretKey, fmt.Sprintf("%s/auth/github/callback", entrypoint)),
		google.New(ggClientId, ggSecretKey, fmt.Sprintf("%s/auth/google/callback", entrypoint)))

	app.Get("/auth/logout", authHandler.Logout)
	app.Get("/auth/user", authHandler.ParseUserFromCookie, authHandler.GetUserInfo)
	app.Get("/auth/:provider", goth_fiber.BeginAuthHandler)
	app.Get("/auth/:provider/callback", authHandler.OauthProviderCallback)

	app.Get("/:token", fileHandler.GetFile)

	err = app.Listen(port)
	if err != nil {
		log.Fatalf("unable to start server on %s: %v", port, err)
	}
}

func getEnv() (string, string, string, time.Duration, string, string) {
	// Required env
	key := os.Getenv("MASTER_KEY")
	storagePath := os.Getenv("STORAGE_PATH")
	if len(key) == 0 || len(storagePath) == 0 {
		log.Fatalln("MASTER_KEY and STORAGE_PATH env must be defined")
	}
	azStorageConnectionString := os.Getenv("AZSTORAGE_CONNECTION_STRING")
	azStorageContainerName := os.Getenv("AZSTORAGE_CONTAINER_NAME")
	if len(azStorageConnectionString) == 0 || len(azStorageContainerName) == 0 {
		log.Fatalln("AZSTORAGE_CONNECTION_STRING and AZSTORAGE_CONTAINER_NAME are required")
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

	return key, storagePath, port, duration, azStorageConnectionString, azStorageContainerName
}

func getDBEnv() (string, string, string, string, string) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_DATABASE")

	return host, port, username, password, dbname
}

func getOauthEnv() (string, string, string, string, string) {
	entrypoint := os.Getenv("ENTRYPOINT")
	ghClientId := os.Getenv("GITHUB_OAUTH_CLIENT_ID")
	ghSecretKey := os.Getenv("GITHUB_OAUTH_SECRET_KEY")
	ggClientId := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	ggSecretKey := os.Getenv("GOOGLE_OAUTH_SECRET_KEY")

	return entrypoint, ghClientId, ghSecretKey, ggClientId, ggSecretKey
}
