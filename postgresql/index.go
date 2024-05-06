package postgresql

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/utilities"
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

	for i := 0; i <= 5; i += 1 {
		db, connectionError := gorm.Open(
			postgres.Open(dsn),
			&gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			},
		)
		if connectionError == nil {
			Database = db
			break
		} else {
			if i == 5 {
				log.Fatal(connectionError)
			}
			log.Printf("Could not connect to PostgreSQL, retrying in %d sec", i+1)
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}

	sqlDB, dbError := Database.DB()
	if dbError != nil {
		log.Fatal(dbError)
	}
	pingError := sqlDB.Ping()
	if pingError != nil {
		log.Fatal(pingError)
	}

	autoMigrationError := Database.AutoMigrate(
		&User{},
		&Password{},
		&UserSecret{},
		&UsedRefreshToken{},
	)
	if autoMigrationError != nil {
		log.Fatal(autoMigrationError)
	}

	log.Println(constants.ACTION_MESSAGES.PGConnected)
}
