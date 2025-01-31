package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"

	apiAccount "go-fiber-auth-2024/apis/account"
	apiAuth "go-fiber-auth-2024/apis/auth"
	apiIndex "go-fiber-auth-2024/apis/index"
	apiUser "go-fiber-auth-2024/apis/user"
	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/middlewares"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/redis"
	"go-fiber-auth-2024/utilities"
)

func main() {
	envSource := utilities.GetEnv(utilities.GetEnvOptions{
		EnvName:    constants.ENV_NAMES.EnvSource,
		IsRequired: true,
	})
	if envSource == constants.ENV_SOURCES.File {
		if envError := godotenv.Load(); envError != nil {
			log.Fatal(envError)
		}
		log.Println(constants.ACTION_MESSAGES.LoadedEnvFile)
	}

	app := fiber.New(fiber.Config{
		AppName:      constants.APP_NAME,
		ErrorHandler: utilities.ErrorHandler,
	})

	app.Use(middlewares.RequestReceivedTimestamp)
	app.Use(favicon.New(favicon.Config{
		File: "./assets/favicon.ico",
		URL:  "/favicon.ico",
	}))
	app.Use(logger.New())

	postgresql.CreateDatabaseConnection()
	redis.CreateDatabaseConnection()

	port := utilities.GetEnv(utilities.GetEnvOptions{
		DefaultValue: constants.PORT,
		EnvName:      constants.ENV_NAMES.Port,
	})

	apiAccount.Initialize(app)
	apiAuth.Initialize(app)
	apiIndex.Initialize(app)
	apiUser.Initialize(app)

	launchError := app.Listen(
		fmt.Sprintf(":%s", port),
		fiber.ListenConfig{
			EnablePrefork:     false,
			OnShutdownError:   utilities.ShutdownError,
			OnShutdownSuccess: utilities.ShutdownSuccess,
		},
	)
	if launchError != nil {
		log.Fatal(launchError)
	}
}
