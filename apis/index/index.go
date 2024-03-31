package index

import "github.com/gofiber/fiber/v3"

func Initialize(server *fiber.App) {
	router := server.Group("/")

	router.Get("/", indexController)
	router.Get("/api", indexController)
}
