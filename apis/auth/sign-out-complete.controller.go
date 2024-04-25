package auth

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/redis"
	"go-fiber-auth-2024/utilities"
)

func signOutCompleteController(context fiber.Ctx) error {
	userId, ok := context.Locals(constants.LOCALS_KEYS.UserId).(int)
	if !ok {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: errors.New(constants.ACTION_MESSAGES.TypeAssertionFailed),
		})
	}

	redisError := redis.Client.Del(
		context.Context(),
		redis.CreateKey(
			constants.REDIS_PREFIXES.SecretHash,
			fmt.Sprint(userId),
		),
	).Err()
	if redisError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: redisError,
		})
	}

	newUserSecret, secretError := utilities.CreateUserSecret(uint(userId))
	if secretError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: secretError,
		})
	}

	tx := postgresql.Database.Begin()
	if tx.Error != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: tx.Error,
		})
	}

	queryError := tx.
		Model(&postgresql.UserSecret{}).
		Where("user_id = ?", userId).Update("secret", newUserSecret).
		Error
	if queryError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}

	queryError = tx.
		Delete(&postgresql.UsedRefreshToken{}, "user_id = ?", userId).
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
