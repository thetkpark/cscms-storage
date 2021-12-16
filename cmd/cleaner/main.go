package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"github.com/thetkpark/cscms-temp-storage/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func main() {
	logger := hclog.Default()

	storagePath := getEnv()
	dbHost, dbPort, dbUsername, dbPassword, dbName := getDBEnv()
	isFailed := false

	// Open data store
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	diskStorageManager, err := service.NewDiskStorageManager(logger, storagePath)
	if err != nil {
		logger.Error("unable to create disk storage manager", err)
		os.Exit(1)
	}

	fileLists, err := diskStorageManager.ListFiles()
	if err != nil {
		logger.Error("unable to list file", err)
		os.Exit(1)
	}

	deletedCount := 0

	for _, fileName := range fileLists {
		var fileInfo model.File
		tx := db.Where(&model.File{ID: fileName}).First(&fileInfo)
		if tx.Error != nil {
			logger.Error(fmt.Sprintf("unable to get file %s info from db", fileName), err)
			isFailed = true
			continue
		}

		// Check if expired_at is in the future
		if fileInfo.ExpiredAt.UTC().After(time.Now().UTC()) {
			continue
		}

		// Delete expired file
		if err := diskStorageManager.DeleteFile(fileName); err != nil {
			// If failed -> continue to delete other file
			isFailed = true
			continue
		}
		deletedCount++
	}

	if isFailed {
		logger.Info("There is an failure")
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Delete %d file", deletedCount))
}

func getEnv() string {
	storagePath := os.Getenv("STORAGE_PATH")
	if len(storagePath) == 0 {
		log.Fatalln("STORAGE_PATH env must be defined")
	}

	return storagePath
}

func getDBEnv() (string, string, string, string, string) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_DATABASE")

	return host, port, username, password, dbname
}
