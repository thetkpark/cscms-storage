package main

type ApplicationEnvironmentVariable struct {
	MasterKey                        string `env:"MASTER_KEY"`
	FileStoragePath                  string `env:"STORAGE_PATH"`
	FileStoreMaxDuration             int    `env:"STORE_DURATION" envDefault:"30"`
	AzureBlobStorageConnectionString string `env:"AZSTORAGE_CONNECTION_STRING"`
	AzureBlobStorageContainerName    string `env:"AZSTORAGE_CONTAINER_NAME"`
	Port                             string `env:"PORT"`
	DB                               DatabaseEnvironmentVariable
	OauthGitHubClientSecret          string `env:"GITHUB_OAUTH_CLIENT_ID"`
	OauthGitHubSecretKey             string `env:"GITHUB_OAUTH_SECRET_KEY"`
	OAuthGoogleClientSecret          string `env:"GOOGLE_OAUTH_CLIENT_ID"`
	OAuthGoogleSecretKey             string `env:"GOOGLE_OAUTH_SECRET_KEY"`
	Entrypoint                       string `env:"ENTRYPOINT"`
}

type DatabaseEnvironmentVariable struct {
	Username     string `env:"DB_USERNAME"`
	Password     string `env:"DB_PASSWORD"`
	Host         string `env:"DB_HOST"`
	Port         string `env:"DB_PORT"`
	DatabaseName string `env:"DB_DATABASE"`
}
