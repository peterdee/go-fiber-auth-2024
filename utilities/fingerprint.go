package utilities

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func Fingerprint(context fiber.Ctx) string {
	return fmt.Sprintf(
		"%s:%s:%s:%s:%s",
		context.Get("user-agent"),
		context.Get("accept-language"),
		context.IP(),
		context.Host(),
		context.Hostname(),
	)
}
