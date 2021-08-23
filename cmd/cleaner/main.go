package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/data/model"
	"github.com/thetkpark/cscms-temp-storage/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func main() {
	logger := hclog.Default()

	storagePath, sqlitePath := getEnv()
	isFailed := false
	// Open data store
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})

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

		// Check if not expire
		if fileInfo.CreatedAt.UTC().Add(time.Hour * 24 * 30).After(time.Now().UTC()) {
			// Not expire -> continue
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

func getEnv() (string, string) {
	storagePath := os.Getenv("STORAGE_PATH")
	dbPath := os.Getenv("SQLITE_PATH")
	if len(storagePath) == 0 || len(dbPath) == 0 {
		log.Fatalln("SQLITE_PATH STORAGE_PATH env must be defined")
	}
	return storagePath, dbPath
}
