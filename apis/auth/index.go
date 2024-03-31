package auth

import "github.com/gofiber/fiber/v3"

func Initialize(server *fiber.App) {
	router := server.Group("/api/auth")

	router.Get("/sign-in", signInController)
	router.Get("/sign-up", signUpController)
}
