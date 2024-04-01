package utilities

import (
	"errors"
	"go-fiber-auth-2024/constants"
	"log"

	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(context fiber.Ctx, err error) error {
	info := constants.RESPONSE_INFO.InternalServerError
	status := fiber.StatusInternalServerError

	var customError *CustomError
	if errors.As(err, &customError) {
		if customError.Info != "" {
			info = customError.Info
		}
		if customError.Status != 0 {
			status = customError.Status
		}
		if customError.Err != nil && customError.Printable {
			log.Println(customError.Err.Error())
		}
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		info = fiberError.Message
		if info == "Internal Server Error" {
			info = constants.RESPONSE_INFO.InternalServerError
		}
		status = fiberError.Code
	}

	return Response(ResponseOptions{
		Context: context,
		Info:    info,
		Status:  status,
	})
}
