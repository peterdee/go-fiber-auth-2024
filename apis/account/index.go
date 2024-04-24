package account

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/middlewares"
)

func Initialize(server *fiber.App) {
	router := server.Group("/api/account")

	router.Get("/", getAccountController, middlewares.Authorization)
}
