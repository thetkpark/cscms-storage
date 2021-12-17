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

	appENVs, err := getAppENVs()
	if err != nil {
		log.Fatalln("Failed to get app ENVs", err)
	}

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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", appENVs.DB.Username, appENVs.DB.Password, appENVs.DB.Host, appENVs.DB.Port, appENVs.DB.DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("unable to open sqlite db", err)
	}
	gormFileDataStore, err := data.NewGormFileDataStore(logger, db, appENVs.FileStoreMaxDuration)
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
	sioEncryptionManager := service.NewSIOEncryptionManager(logger, appENVs.MasterKey)
	diskStorageManager, err := service.NewDiskStorageManager(logger, appENVs.FileStoragePath)
	if err != nil {
		log.Fatalln("unable to create disk storage manager")
	}
	imageStorageManager, err := service.NewAzureImageStorageManager(logger, appENVs.AzureBlobStorageConnectionString, appENVs.AzureBlobStorageContainerName)
	if err != nil {
		log.Fatalln("unable to azure image storage manager")
	}
	jwtManager := service.NewJwtManager(os.Getenv("JWT_SECRET"))

	// Create handlers
	fileHandler := handlers.NewFileRoutesHandler(logger, sioEncryptionManager, gormFileDataStore, diskStorageManager, appENVs.FileStoreMaxDuration)
	imageHandler := handlers.NewImageRouteHandler(logger, gormImageDataStore, imageStorageManager)
	authHandler := handlers.NewAuthRouteHandler(logger, gormUserDataStore, jwtManager, appENVs.Entrypoint)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://storage.cscms.me, http://localhost:5050",
		AllowMethods:     "GET POST PATCH DELETE",
		AllowCredentials: true,
	}))

	app.Get("/api/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":   "pong",
			"timestamp": time.Now(),
		})
	})
	app.Post("/api/file", authHandler.ParseUserFromCookie, fileHandler.UploadFile)
	//app.Get("/api/file/lists", authHandler.ParseUserFromCookie, authHandler.AuthenticatedOnly)

	app.Post("/api/image", imageHandler.UploadImage)
	//app.Get("/api/image/lists", authHandler.ParseUserFromCookie, authHandler.AuthenticatedOnly)

	app.Static("/", "./client/build")
	app.Static("/404", "./client/build")

	// User Authentication with Oauth
	goth.UseProviders(
		github.New(appENVs.OauthGitHub.ClientSecret, appENVs.OauthGitHub.SecretKey, fmt.Sprintf("%s/auth/github/callback", appENVs.Entrypoint)),
		google.New(appENVs.OAuthGoogle.ClientSecret, appENVs.OAuthGoogle.SecretKey, fmt.Sprintf("%s/auth/google/callback", appENVs.Entrypoint)))

	app.Get("/auth/logout", authHandler.Logout)
	app.Get("/auth/user", authHandler.ParseUserFromCookie, authHandler.GetUserInfo)
	app.Get("/auth/:provider", goth_fiber.BeginAuthHandler)
	app.Get("/auth/:provider/callback", authHandler.OauthProviderCallback)

	app.Get("/:token", fileHandler.GetFile)

	err = app.Listen(appENVs.Port)
	if err != nil {
		log.Fatalf("unable to start server on %s: %v", appENVs.Port, err)
	}
}
