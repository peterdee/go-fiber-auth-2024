package middlewares

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/julyskies/gohelpers"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/redis"
	"go-fiber-auth-2024/utilities"
)

func Authorization(context fiber.Ctx) error {
	accessTokenRaw := context.Get("authorization")
	if accessTokenRaw == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingToken,
			Status: fiber.StatusUnauthorized,
		})
	}
	accessToken := strings.Trim(accessTokenRaw, " ")
	if accessToken == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	claims, decodeError := utilities.DecodeToken(accessToken)
	if decodeError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.InvlaidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	issuedAtSeconds := claims.Issued.Time().UnixNano() / int64(time.Millisecond)
	tokenPairId := claims.ID
	userIdString := claims.Subject

	if issuedAtSeconds == 0 || tokenPairId == "" || userIdString == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.InvlaidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	accessTokenExpirationString := utilities.GetEnv(utilities.GetEnvOptions{
		DefaultValue: fmt.Sprint(constants.TOKENS.DefaultAccessTokenExpirationSeconds),
		EnvName:      constants.ENV_NAMES.AccessTokenExpirationSeconds,
	})
	accessTokenExpiration, convertError := strconv.Atoi(accessTokenExpirationString)
	if convertError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: convertError,
		})
	}
	if issuedAtSeconds+int64(accessTokenExpiration) < gohelpers.MakeTimestamp() {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.AccessTokenExpired,
			Status: fiber.StatusUnauthorized,
		})
	}

	fingerprint, fingerprintError := utilities.Fingerprint(context)
	if fingerprintError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: fingerprintError,
		})
	}

	// TODO: get user secret hash & user password hash for token secret
	userId, convertError := strconv.Atoi(userIdString)
	if convertError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: convertError,
		})
	}

	userSecretHash, redisError := redis.Client.Get(context.Context(), redis.CreateKey(
		constants.REDIS_PREFIXES.SecretHash,
		userIdString,
	)).Result()
	if redisError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: redisError,
		})
	}
	if userSecretHash == "" {
		var userSecretRecord postgresql.UserSecret
		queryError := postgresql.
			Database.
			Where(&postgresql.UserSecret{UserID: uint(userId)}).
			First(&userSecretRecord).
			Error
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
		setError := redis.Client.Set(
			context.Context(),
			redis.CreateKey(
				constants.REDIS_PREFIXES.SecretHash,
				userIdString,
			),
			userSecretRecord.Secret,
			time.Hour*4,
		).Err()
		if setError != nil {
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: setError,
			})
		}
		userSecretHash = userSecretRecord.Secret
	}

	userPasswordHash, redisError := redis.Client.Get(
		context.Context(),
		redis.CreateKey(
			constants.REDIS_PREFIXES.PasswordHash,
			userIdString,
		),
	).Result()
	if redisError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: redisError,
		})
	}
	if userPasswordHash == "" {
		var userPasswordRecord postgresql.Password
		queryError := postgresql.
			Database.
			Where(&postgresql.Password{UserID: uint(userId)}).
			First(&userPasswordRecord).
			Error
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
		setError := redis.Client.Set(
			context.Context(),
			redis.CreateKey(
				constants.REDIS_PREFIXES.SecretHash,
				userIdString,
			),
			userPasswordRecord.Hash,
			time.Hour*4,
		).Err()
		if setError != nil {
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: setError,
			})
		}
		userPasswordHash = userPasswordRecord.Hash
	}

	tokenSecret, tokenSecretError := utilities.CreateTokenSecret(
		userSecretHash,
		userPasswordHash,
		utilities.GetEnv(utilities.GetEnvOptions{
			DefaultValue: constants.TOKENS.DefaultAccessTokenCommonSecret,
			EnvName:      constants.ENV_NAMES.AccessTokenCommonSecret,
		}),
		fingerprint,
	)
	if tokenSecretError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: tokenSecretError,
		})
	}

	tokenIsValid := utilities.VerifyToken(accessToken, tokenSecret)
	if !tokenIsValid {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.InvlaidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	context.Locals(constants.LOCALS_KEYS.TokenPairId, tokenPairId)
	context.Locals(constants.LOCALS_KEYS.UserId, userId)

	return context.Next()
}
