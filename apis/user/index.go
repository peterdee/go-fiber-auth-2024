package user

import "github.com/gofiber/fiber/v3"

func Initialize(server *fiber.App) {
	router := server.Group("/api/user")

	router.Get("/:id", getUser)
}
