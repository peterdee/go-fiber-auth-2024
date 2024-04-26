package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/redis"
	"go-fiber-auth-2024/utilities"
)

func signOutFromAllDevices(userId int, tx *gorm.DB, context fiber.Ctx) error {
	newUserSecretHash, internalError := utilities.CreateUserSecret(uint(userId))
	if internalError != nil {
		return internalError
	}

	queryError := tx.
		Model(&postgresql.UserSecret{}).
		Where("user_id = ?", userId).
		Update("secret", newUserSecretHash).
		Error
	if queryError != nil {
		return queryError
	}

	queryError = tx.
		Delete(&postgresql.UsedRefreshToken{}, "user_id = ?", userId).
		Error
	if queryError != nil {
		return queryError
	}

	redisError := redis.Client.Del(
		context.Context(),
		redis.CreateKey(
			constants.REDIS_PREFIXES.SecretHash,
			fmt.Sprint(userId),
		),
	).Err()
	if redisError != nil {
		return redisError
	}

	return nil
}
