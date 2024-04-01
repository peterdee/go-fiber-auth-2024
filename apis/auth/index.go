package auth

import "github.com/gofiber/fiber/v3"

func Initialize(server *fiber.App) {
	router := server.Group("/api/auth")

	router.Post("/sign-in", signInController)
	router.Post("/sign-up", signUpController)
}
