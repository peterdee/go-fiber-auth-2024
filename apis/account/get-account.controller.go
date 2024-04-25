package account

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/utilities"
)

func getAccountController(context fiber.Ctx) error {
	userId, ok := context.Locals(constants.LOCALS_KEYS.UserId).(int)
	if !ok {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: errors.New(constants.ACTION_MESSAGES.TypeAssertionFailed),
		})
	}

	var account postgresql.User
	queryError := postgresql.Database.Where("id = ?", userId).First(&account).Error
	if queryError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}
	return utilities.Response(utilities.ResponseOptions{
		Context: context,
		Data: fiber.Map{
			"account": account,
		},
	})
}
