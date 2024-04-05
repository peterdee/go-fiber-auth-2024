package postgresql

import (
	"fmt"
	"go-fiber-auth-2024/constants"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func getEnvValue(envName string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Fatalf("%s: %s", constants.ACTION_MESSAGES.PGCredentialsError, envName)
	}
	return value
}

func CreateDatabaseConnection() {
	database := getEnvValue(constants.ENV_NAMES.PGDatabase)
	host := getEnvValue(constants.ENV_NAMES.PGHost)
	password := getEnvValue(constants.ENV_NAMES.PGPassword)
	port := getEnvValue(constants.ENV_NAMES.PGPort)
	username := getEnvValue(constants.ENV_NAMES.PGUsername)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		username,
		password,
		database,
		port,
	)

	db, connectionError := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if connectionError != nil {
		log.Fatal(connectionError)
	}

	sqlDB, dbError := db.DB()
	if dbError != nil {
		log.Fatal(dbError)
	}
	pingError := sqlDB.Ping()
	if pingError != nil {
		log.Fatal(pingError)
	}

	autoMigrationError := db.AutoMigrate(
		&User{},
		&Password{},
		&UserSecret{},
		&UsedRefreshToken{},
	)
	if autoMigrationError != nil {
		log.Fatal(autoMigrationError)
	}

	Database = db

	log.Println(constants.ACTION_MESSAGES.PGConnected)
}
