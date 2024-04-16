package account

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/utilities"
)

func changePasswordController(context fiber.Ctx) error {
	if len(context.Body()) == 0 {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	payload := new(ChangePasswordPayload)
	if bindingError := context.Bind().Body(payload); bindingError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: bindingError,
		})
	}
	if payload.NewPassword == "" || payload.OldPassword == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	userId, ok := context.Locals(constants.LOCALS_KEYS.UserId).(int)
	if !ok {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: errors.New(constants.ACTION_MESSAGES.TypeAssertionFailed),
		})
	}

	newPassword := strings.Trim(payload.NewPassword, " ")
	oldPassword := strings.Trim(payload.OldPassword, " ")

	var userPasswordRecord postgresql.Password
	queryError := postgresql.Database.Where("user_id = ?", userId).First(&userPasswordRecord).Error
	if queryError != nil {
		if queryError.Error() == "record not found" {
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Info:   constants.RESPONSE_INFO.Unauthorized,
				Status: fiber.StatusUnauthorized,
			})
		}
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}

	isValid, internalError := utilities.ComparePasswordAndHash(
		oldPassword,
		userPasswordRecord.Hash,
	)
	if internalError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: internalError,
		})
	}
	if !isValid {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.OldPasswordIsInvalid,
			Status: fiber.StatusBadRequest,
		})
	}

	newPasswordHash, internalError := utilities.CreateHash(newPassword)
	if internalError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: internalError,
		})
	}

	fmt.Println(newPasswordHash)

	// TODO: update password record, delete password hash from Redis, wrap into transaction

	return utilities.Response(utilities.ResponseOptions{Context: context})
}
