package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/thetkpark/cscms-temp-storage/data"
	"github.com/thetkpark/cscms-temp-storage/service"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func main() {

	// Get ENV
	appENVs := ApplicationEnvironmentVariable{}
	if err := env.Parse(&appENVs, env.Options{RequiredIfNoDef: true}); err != nil {
		log.Fatalf("Unable to get env: %v", err.Error())
	}

	zapLogger, _ := zap.NewProduction()
	if appENVs.Env == "development" {
		zapLogger, _ = zap.NewDevelopment()
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	// Open data store
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", appENVs.DB.Username, appENVs.DB.Password, appENVs.DB.Host, appENVs.DB.Port, appENVs.DB.DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Errorw("unable to open connection to db", "error", err.Error())
	}

	// Create disk storage manager
	diskStorageManager, err := service.NewDiskStorageManager(logger, appENVs.FileStoragePath)
	if err != nil {
		logger.Errorw("unable to create disk storage manager", "error", err.Error())
		os.Exit(1)
	}
	// Create file data store
	fileDataStore, err := data.NewGormFileDataStore(db, time.Duration(appENVs.FileStoreMaxDuration)*time.Hour*24)
	if err != nil {
		logger.Errorw("unable to create file data store", "error", err.Error())
	}

	fileLists, err := diskStorageManager.ListFiles()
	if err != nil {
		logger.Errorw("unable to list file", "error", err.Error())
		os.Exit(1)
	}

	deletedCount := 0
	isError := false

	for _, fileName := range fileLists {
		fileInfo, err := fileDataStore.FindByID(fileName)
		if err != nil {
			isError = true
			logger.Errorw("Unable to query by file id", "error", err.Error())
			continue
		}

		// Check if fileInfo is existed in db and expired_at is in the future
		if fileInfo != nil && fileInfo.ExpiredAt.UTC().After(time.Now().UTC()) {
			continue
		}

		// Delete expired file
		if err := diskStorageManager.DeleteFile(fileName); err != nil {
			// If failed -> continue to delete other file
			isError = true
			continue
		}
		deletedCount++
	}

	if isError {
		logger.Info("There is an failure")
	}

	logger.Info(fmt.Sprintf("Delete %d file", deletedCount))
}

type ApplicationEnvironmentVariable struct {
	FileStoragePath      string `env:"STORAGE_PATH"`
	FileStoreMaxDuration int    `env:"STORE_DURATION" envDefault:"30"`
	Env                  string `env:"ENV" envDefault:"development"`
	DB                   DatabaseEnvironmentVariable
}

type DatabaseEnvironmentVariable struct {
	Username     string `env:"DB_USERNAME"`
	Password     string `env:"DB_PASSWORD"`
	Host         string `env:"DB_HOST"`
	Port         string `env:"DB_PORT"`
	DatabaseName string `env:"DB_DATABASE"`
}
