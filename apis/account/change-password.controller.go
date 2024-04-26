package account

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/redis"
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

	tx := postgresql.Database.Begin()
	if tx.Error != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: tx.Error,
		})
	}

	var userPasswordRecord postgresql.Password
	queryError := tx.
		Where("user_id = ?", userId).
		First(&userPasswordRecord).
		Error
	if queryError != nil {
		tx.Rollback()
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
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: internalError,
		})
	}
	if !isValid {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.OldPasswordIsInvalid,
			Status: fiber.StatusBadRequest,
		})
	}

	newPasswordHash, internalError := utilities.CreateHash(newPassword)
	if internalError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: internalError,
		})
	}

	redisError := redis.Client.Del(
		context.Context(),
		redis.CreateKey(
			constants.REDIS_PREFIXES.PasswordHash,
			fmt.Sprint(userId),
		),
	).Err()
	if redisError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: redisError,
		})
	}

	queryError = tx.
		Model(&postgresql.Password{}).
		Where("user_id = ?", userId).Update("hash", newPasswordHash).
		Error
	if queryError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}

	commitError := tx.Commit().Error
	if commitError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: commitError,
		})
	}

	return utilities.Response(utilities.ResponseOptions{Context: context})
}
