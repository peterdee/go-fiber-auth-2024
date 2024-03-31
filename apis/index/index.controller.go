package index

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/utilities"
)

func indexController(context fiber.Ctx) error {
	return utilities.Response(utilities.ResponseOptions{Context: context})
}
