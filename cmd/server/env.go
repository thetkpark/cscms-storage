package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type ApplicationEnvironmentVariable struct {
	MasterKey                        string
	FileStoragePath                  string
	FileStoreMaxDuration             string
	AzureBlobStorageConnectionString string
	AzureBlobStorageContainerName    string
	Port                             string
	DB                               DatabaseEnvironmentVariable
	OauthGitHub                      OauthEnvironmentVariable
	OAuthGoogle                      OauthEnvironmentVariable
}

type DatabaseEnvironmentVariable struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
}

type OauthEnvironmentVariable struct {
	ClientSecret string
	SecretKey    string
	CallbackURL  string
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
