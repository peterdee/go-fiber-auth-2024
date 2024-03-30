package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/julyskies/gohelpers"

	"go-fiber-auth-2024/constants"
)

func RequestReceivedTimestamp(context fiber.Ctx) error {
	context.Locals(
		constants.LOCALS_KEYS.RequestReceivedTimestamp,
		gohelpers.MakeTimestamp(),
	)
	return context.Next()
}
