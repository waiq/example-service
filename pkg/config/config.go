package config

import (
	"fmt"

	"github.com/nicholasjackson/env"
)

var ApplicationPort = env.Int("APP_PORT", false, 8080, "Port fot this application server")

var DatabasePort = env.Int("DB_PORT", false, 5432, "Database port")
var DatabaseName = env.String("DB_NAME", false, "books", "Database name")
var DatabaseUser = env.String("DB_USER", false, "booksuser", "Database host")
var DatabasePassword = env.String("DB_PASSWORD", true, "localhost", "Database host")
var DatabaseHost = env.String("DB_HOST", true, "localhost", "Database host")

func InitConfig() {
	err := env.Parse()
	if err != nil {
		panic(err)
	}
}

func GetPostgresDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Etc/UTC",
		*DatabaseHost,
		*DatabaseUser,
		*DatabasePassword,
		*DatabaseName,
		*DatabasePort,
	)
}
