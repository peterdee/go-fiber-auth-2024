package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/julyskies/gohelpers"
	"gorm.io/gorm"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/redis"
	"go-fiber-auth-2024/utilities"
)

func createTokens(
	userId uint,
	tx *gorm.DB,
	context fiber.Ctx,
) (string, string, error) {
	fingerprint := utilities.Fingerprint(context)
	tokenPairId := gohelpers.RandomString(constants.TOKEN_PAIR_ID_LENGTH)

	var userPasswordRecord postgresql.Password
	queryError := tx.
		Where(&postgresql.Password{UserID: userId}).
		First(&userPasswordRecord).
		Error
	if queryError != nil {
		return "", "", queryError
	}

	var userSecretRecord postgresql.UserSecret
	queryError = tx.
		Where(&postgresql.UserSecret{UserID: userId}).
		First(&userSecretRecord).
		Error
	if queryError != nil {
		return "", "", queryError
	}

	accessTokenSecret := utilities.CreateTokenSecret(
		userSecretRecord.Secret,
		userPasswordRecord.Hash,
		utilities.GetEnv(utilities.GetEnvOptions{
			DefaultValue: constants.TOKENS.DefaultAccessTokenCommonSecret,
			EnvName:      constants.ENV_NAMES.AccessTokenCommonSecret,
		}),
		fingerprint,
	)
	accessToken, tokenError := utilities.CreateToken(
		fmt.Sprint(userId),
		accessTokenSecret,
		tokenPairId,
	)
	if tokenError != nil {
		return "", "", tokenError
	}

	refreshTokenSecret := utilities.CreateTokenSecret(
		userSecretRecord.Secret,
		userPasswordRecord.Hash,
		utilities.GetEnv(utilities.GetEnvOptions{
			DefaultValue: constants.TOKENS.DefaultRefreshTokenCommonSecret,
			EnvName:      constants.ENV_NAMES.RefreshTokenCommonSecret,
		}),
		fingerprint,
	)
	refreshToken, tokenError := utilities.CreateToken(
		fmt.Sprint(userId),
		refreshTokenSecret,
		tokenPairId,
	)
	if tokenError != nil {
		return "", "", tokenError
	}

	return accessToken, refreshToken, nil
}

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
