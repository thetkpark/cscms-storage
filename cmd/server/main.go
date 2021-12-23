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

	"github.com/arsmn/fiber-swagger/v2"
	"github.com/caarlos0/env/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	//"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
	_ "github.com/thetkpark/cscms-temp-storage/cmd/server/docs"
)

// @title CSCMS Storage
// @version 1.0
// @description This is documentation for CSCMS Storage API

func main() {
	logger := hclog.Default()

	appENVs := ApplicationEnvironmentVariable{}
	if err := env.Parse(&appENVs, env.Options{RequiredIfNoDef: true}); err != nil {
		log.Fatalln("Failed to get app ENVs: ", err)
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
	gormFileDataStore, err := data.NewGormFileDataStore(logger, db, time.Duration(appENVs.FileStoreMaxDuration)*time.Hour*24)
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
	fileHandler := handlers.NewFileRoutesHandler(logger, sioEncryptionManager, gormFileDataStore, diskStorageManager, time.Duration(appENVs.FileStoreMaxDuration)*time.Hour*24)
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

	apiPath := app.Group("/api", authHandler.ParseUserFromCookie)

	filePath := apiPath.Group("/file")
	filePath.Post("/", fileHandler.UploadFile)
	filePath.Get("/", authHandler.AuthenticatedOnly, fileHandler.GetOwnFiles)
	filePath.Patch("/:fileID", authHandler.AuthenticatedOnly, fileHandler.IsOwnFile, fileHandler.EditToken)
	filePath.Delete("/:fileID", authHandler.AuthenticatedOnly, fileHandler.IsOwnFile, fileHandler.DeleteFile)

	imagePath := apiPath.Group("/image")
	imagePath.Post("/", imageHandler.UploadImage)
	imagePath.Get("/", authHandler.AuthenticatedOnly, imageHandler.GetOwnImages)
	imagePath.Delete("/:imageID", authHandler.AuthenticatedOnly, imageHandler.IsOwnImage, imageHandler.DeleteImage)

	// User Authentication with Oauth
	goth.UseProviders(
		github.New(appENVs.OauthGitHubClientSecret, appENVs.OauthGitHubSecretKey, fmt.Sprintf("%s/auth/github/callback", appENVs.Entrypoint)),
		google.New(appENVs.OAuthGoogleClientSecret, appENVs.OAuthGoogleSecretKey, fmt.Sprintf("%s/auth/google/callback", appENVs.Entrypoint)))

	authPath := app.Group("/auth")
	authPath.Get("/logout", authHandler.Logout)
	authPath.Get("/user", authHandler.ParseUserFromCookie, authHandler.AuthenticatedOnly, authHandler.GetUserInfo)
	authPath.Get("/:provider", goth_fiber.BeginAuthHandler)
	authPath.Get("/:provider/callback", authHandler.OauthProviderCallback)

	// Other routes
	app.Static("/", "./client/build")
	app.Static("/404", "./client/build")
	app.Get("/swagger/*", swagger.Handler)
	app.Get("/:token", fileHandler.GetFile)

	err = app.Listen(fmt.Sprintf(":%s", appENVs.Port))
	if err != nil {
		log.Fatalf("unable to start server on %s: %v", appENVs.Port, err)
	}
}
