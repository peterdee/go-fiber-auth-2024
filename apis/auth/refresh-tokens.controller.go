package auth

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/utilities"
)

func refreshTokensController(context fiber.Ctx) error {
	return utilities.Response(utilities.ResponseOptions{Context: context})
}
