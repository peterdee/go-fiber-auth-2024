package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"

	apiAuth "go-fiber-auth-2024/apis/auth"
	apiIndex "go-fiber-auth-2024/apis/index"
	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/middlewares"
	"go-fiber-auth-2024/utilities"
)

func main() {
	envSource := os.Getenv(constants.ENV_NAMES.EnvSource)
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

	port := constants.PORT
	if envPort := os.Getenv(constants.ENV_NAMES.Port); envPort != "" {
		port = envPort
	}

	apiAuth.Initialize(app)
	apiIndex.Initialize(app)

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
