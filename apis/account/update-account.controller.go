package account

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/utilities"
)

func updateAccountController(context fiber.Ctx) error {
	if len(context.Body()) == 0 {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	payload := new(UpdateAccountPayload)
	if bindingError := context.Bind().Body(payload); bindingError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: bindingError,
		})
	}

	userId, ok := context.Locals(constants.LOCALS_KEYS.UserId).(int)
	if !ok {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: errors.New(constants.ACTION_MESSAGES.TypeAssertionFailed),
		})
	}

	queryError := postgresql.Database.
		Model(&postgresql.User{}).
		Where("id = ?", userId).
		Updates(&payload).
		Error
	if queryError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}

	return utilities.Response(utilities.ResponseOptions{
		Context: context,
	})
}
