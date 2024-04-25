package auth

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/middlewares"
)

func Initialize(server *fiber.App) {
	router := server.Group("/api/auth")

	router.Post("/refresh", refreshTokensController)
	router.Post("/sign-in", signInController)
	router.Post("/sign-out", signOutController, middlewares.Authorization)
	router.Get("/sign-out/complete", signOutCompleteController, middlewares.Authorization)
	router.Post("/sign-up", signUpController)
}
