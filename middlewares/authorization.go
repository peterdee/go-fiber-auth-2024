package middlewares

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"

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
			Info:   constants.RESPONSE_INFO.InvalidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	issuedAtSeconds := claims.Issued.Time().Unix()
	tokenPairId := claims.ID
	userIdString := claims.Subject

	if issuedAtSeconds == 0 || tokenPairId == "" || userIdString == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.InvalidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	isExpired, internalError := utilities.CheckTokenExpiration(
		issuedAtSeconds,
		utilities.TOKEN_TYPE_ACCESS,
	)
	if internalError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: internalError,
		})
	}
	if isExpired {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.AccessTokenExpired,
			Status: fiber.StatusUnauthorized,
		})
	}

	accessTokenExpirationSecondsString := utilities.GetEnv(utilities.GetEnvOptions{
		DefaultValue: fmt.Sprint(constants.TOKENS.DefaultAccessTokenExpirationSeconds),
		EnvName:      constants.ENV_NAMES.AccessTokenExpirationSeconds,
	})
	accessTokenExpirationSeconds, convertError := strconv.Atoi(
		accessTokenExpirationSecondsString,
	)
	if convertError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: convertError,
		})
	}
	tokenPairIdKey := redis.CreateKey(
		constants.REDIS_PREFIXES.BlacklistedTokenPair,
		fmt.Sprintf("%s-%s", userIdString, claims.ID),
	)
	blacklistedTokenPairId, redisError := redis.Client.Get(
		context.Context(),
		tokenPairIdKey,
	).Result()
	if redisError != nil && redisError.Error() != "redis: nil" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: redisError,
		})
	}
	if blacklistedTokenPairId != "" {
		expireError := redis.Client.Expire(
			context.Context(),
			tokenPairIdKey,
			time.Duration(accessTokenExpirationSeconds)*time.Second,
		).Err()
		if expireError != nil {
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: expireError,
			})
		}
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.AccessTokenBlacklisted,
			Status: fiber.StatusUnauthorized,
		})
	}

	fingerprint := utilities.Fingerprint(context)

	userId, convertError := strconv.Atoi(userIdString)
	if convertError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: convertError,
		})
	}

	userSecretHashKey := redis.CreateKey(
		constants.REDIS_PREFIXES.PasswordHash,
		userIdString,
	)
	userSecretHash, redisError := redis.Client.Get(
		context.Context(),
		userSecretHashKey,
	).Result()
	if redisError != nil && redisError.Error() != "redis: nil" {
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
	} else {
		expireError := redis.Client.Expire(
			context.Context(),
			userSecretHashKey,
			time.Duration(accessTokenExpirationSeconds)*time.Second,
		).Err()
		if expireError != nil {
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: expireError,
			})
		}
	}

	userPasswordHashKey := redis.CreateKey(
		constants.REDIS_PREFIXES.PasswordHash,
		userIdString,
	)
	userPasswordHash, redisError := redis.Client.Get(
		context.Context(),
		userPasswordHashKey,
	).Result()
	if redisError != nil && redisError.Error() != "redis: nil" {
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
	} else {
		expireError := redis.Client.Expire(
			context.Context(),
			userPasswordHashKey,
			time.Duration(accessTokenExpirationSeconds)*time.Second,
		).Err()
		if expireError != nil {
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: expireError,
			})
		}
	}

	tokenSecret := utilities.CreateTokenSecret(
		userSecretHash,
		userPasswordHash,
		utilities.GetEnv(utilities.GetEnvOptions{
			DefaultValue: constants.TOKENS.DefaultAccessTokenCommonSecret,
			EnvName:      constants.ENV_NAMES.AccessTokenCommonSecret,
		}),
		fingerprint,
	)

	tokenIsValid := utilities.VerifyToken(accessToken, tokenSecret)
	if !tokenIsValid {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.InvalidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	context.Locals(constants.LOCALS_KEYS.TokenPairId, tokenPairId)
	context.Locals(constants.LOCALS_KEYS.UserId, userId)

	return context.Next()
}
