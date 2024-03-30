package utilities

import (
	"github.com/gofiber/fiber/v3"
	"github.com/julyskies/gohelpers"

	"go-fiber-auth-2024/constants"
)

func Response(options ResponseOptions) error {
	info := constants.RESPONSE_INFO.Ok
	if options.Info != "" {
		info = options.Info
	}
	status := fiber.StatusOK
	if options.Status != 0 {
		status = options.Status
	}

	now := gohelpers.MakeTimestamp()
	response := fiber.Map{
		"datetime":   now,
		"info":       info,
		"requestURL": options.Context.OriginalURL() + " [" + options.Context.Method() + "]",
		"status":     status,
	}
	if options.Data != nil {
		response["data"] = options.Data
	}

	requestReceivedTimestamp := options.Context.Locals(
		constants.LOCALS_KEYS.RequestReceivedTimestamp,
	)
	if val, ok := requestReceivedTimestamp.(int64); ok {
		response["processedMS"] = now - val
	}

	return options.Context.Status(options.Status).JSON(response)
}
