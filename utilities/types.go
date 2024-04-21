package utilities

import "github.com/gofiber/fiber/v3"

type GetEnvOptions struct {
	DefaultValue string
	EnvName      string
	IsRequired   bool
}

type ResponseOptions struct {
	Context fiber.Ctx
	Data    fiber.Map
	Info    string
	Status  int
}
