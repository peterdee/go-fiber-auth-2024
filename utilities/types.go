package utilities

import "github.com/gofiber/fiber/v3"

type ResponseOptions struct {
	Context fiber.Ctx
	Data    fiber.Map
	Info    string
	Status  int
}
