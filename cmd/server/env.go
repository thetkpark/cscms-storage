package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type ApplicationEnvironmentVariable struct {
	MasterKey                        string
	FileStoragePath                  string
	FileStoreMaxDuration             time.Duration
	AzureBlobStorageConnectionString string
	AzureBlobStorageContainerName    string
	Port                             string
	DB                               DatabaseEnvironmentVariable
	OauthGitHub                      OauthEnvironmentVariable
	OAuthGoogle                      OauthEnvironmentVariable
	Entrypoint                       string
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
}

func getEnv() (*ApplicationEnvironmentVariable, error) {
	requireEnv := []string{"MASTER_KEY", "STORAGE_PATH", "AZSTORAGE_CONNECTION_STRING", "AZSTORAGE_CONTAINER_NAME"}
	missingENVs := make([]string, 0, len(requireEnv))

	appEnv := &ApplicationEnvironmentVariable{
		MasterKey:                        getAndCheckRequireENV("MASTER_KEY", &missingENVs),
		FileStoragePath:                  getAndCheckRequireENV("STORAGE_PATH", &missingENVs),
		FileStoreMaxDuration:             time.Hour * 24 * 30,
		AzureBlobStorageConnectionString: getAndCheckRequireENV("AZSTORAGE_CONNECTION_STRING", &missingENVs),
		AzureBlobStorageContainerName:    getAndCheckRequireENV("AZSTORAGE_CONTAINER_NAME", &missingENVs),
		Port:                             ":5000",
		Entrypoint:                       getAndCheckRequireENV("ENTRYPOINT", &missingENVs),
		DB: DatabaseEnvironmentVariable{
			Username:     getAndCheckRequireENV("DB_USERNAME", &missingENVs),
			Password:     getAndCheckRequireENV("DB_PASSWORD", &missingENVs),
			Host:         getAndCheckRequireENV("DB_PORT", &missingENVs),
			Port:         getAndCheckRequireENV("DB_HOST", &missingENVs),
			DatabaseName: getAndCheckRequireENV("DB_DATABASE", &missingENVs),
		},
		OauthGitHub: OauthEnvironmentVariable{
			ClientSecret: getAndCheckRequireENV("GITHUB_OAUTH_CLIENT_ID", &missingENVs),
			SecretKey:    getAndCheckRequireENV("GITHUB_OAUTH_SECRET_KEY", &missingENVs),
		},
		OAuthGoogle: OauthEnvironmentVariable{
			ClientSecret: getAndCheckRequireENV("GOOGLE_OAUTH_CLIENT_ID", &missingENVs),
			SecretKey:    getAndCheckRequireENV("GOOGLE_OAUTH_SECRET_KEY", &missingENVs),
		},
	}

	// Check missing require ENV
	if len(missingENVs) != 0 {
		errorString := "Missing ENV: "
		for _, envKey := range missingENVs {
			errorString += envKey + " "
		}
		return nil, fmt.Errorf(errorString)
	}

	// Optional env
	port := os.Getenv("PORT")
	if len(port) != 0 {
		appEnv.Port = fmt.Sprintf(":%s", port)
	}

	maxStoreDuration := os.Getenv("STORE_DURATION") // in days
	if len(maxStoreDuration) != 0 {
		date, err := strconv.Atoi(maxStoreDuration)
		if err != nil {
			return nil, fmt.Errorf("STORE_DURATION is not a valid number")

		}
		appEnv.FileStoreMaxDuration = time.Hour * 24 * time.Duration(date)
	}

	return appEnv, nil
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

func getAndCheckRequireENV(envName string, missingEnv *[]string) string {
	value := os.Getenv(envName)
	if len(value) == 0 {
		*missingEnv = append(*missingEnv, envName)
	}
	return value
}
