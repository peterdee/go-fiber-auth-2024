package auth

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/middlewares"
)

func Initialize(server *fiber.App) {
	router := server.Group("/api/auth")

	router.Post("/sign-in", signInController)
	router.Post("/sign-out", middlewares.Authorization, signOutController)
	router.Post("/sign-up", signUpController)
}
