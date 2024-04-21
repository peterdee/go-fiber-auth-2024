package postgresql

import (
	"fmt"
	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/utilities"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func CreateDatabaseConnection() {
	database := utilities.GetEnv(utilities.GetEnvOptions{
		EnvName:    constants.ENV_NAMES.PGDatabase,
		IsRequired: true,
	})
	host := utilities.GetEnv(utilities.GetEnvOptions{
		EnvName:    constants.ENV_NAMES.PGHost,
		IsRequired: true,
	})
	password := utilities.GetEnv(utilities.GetEnvOptions{
		EnvName:    constants.ENV_NAMES.PGPassword,
		IsRequired: true,
	})
	port := utilities.GetEnv(utilities.GetEnvOptions{
		EnvName:    constants.ENV_NAMES.PGPort,
		IsRequired: true,
	})
	username := utilities.GetEnv(utilities.GetEnvOptions{
		EnvName:    constants.ENV_NAMES.PGUsername,
		IsRequired: true,
	})

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
