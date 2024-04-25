package auth

import (
	"errors"
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

func signOutController(context fiber.Ctx) error {
	if len(context.Body()) == 0 {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	payload := new(SignOutPayload)
	if bindingError := context.Bind().Body(payload); bindingError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: bindingError,
		})
	}
	if payload.RefreshToken == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	accessTokenPairId := context.Locals(constants.LOCALS_KEYS.TokenPairId)
	userId, ok := context.Locals(constants.LOCALS_KEYS.UserId).(int)
	if !ok {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: errors.New(constants.ACTION_MESSAGES.TypeAssertionFailed),
		})
	}

	refreshToken := strings.Trim(payload.RefreshToken, " ")
	refreshTokenClaims, decodeError := utilities.DecodeToken(refreshToken)
	if decodeError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: decodeError,
		})
	}
	refreshTokenPairId := refreshTokenClaims.ID

	tx := postgresql.Database.Begin()
	if tx.Error != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: tx.Error,
		})
	}

	var refreshTokenRecord postgresql.UsedRefreshToken
	queryError := tx.
		Where(&postgresql.UsedRefreshToken{Token: refreshToken, UserID: uint(userId)}).
		First(&refreshTokenRecord).
		Error
	if queryError != nil && queryError.Error() != "record not found" {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}
	if (refreshTokenRecord.ID != 0) ||
		(accessTokenPairId != refreshTokenPairId) {
		secret, secretError := utilities.CreateUserSecret(uint(userId))
		if secretError != nil {
			tx.Rollback()
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: secretError,
			})
		}
		queryError := tx.
			Model(&postgresql.UserSecret{}).
			Where("user_id = ?", userId).Update("secret", secret).
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

		delError := redis.Client.Del(
			context.Context(),
			redis.CreateKey(
				constants.REDIS_PREFIXES.SecretHash,
				fmt.Sprint(userId),
			),
		).Err()
		if delError != nil {
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: delError,
			})
		}

		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.Unauthorized,
			Status: fiber.StatusUnauthorized,
		})
	}

	issuedAtSeconds := refreshTokenClaims.Issued.Time().Unix()
	refreshTokenExpirationSecondsString := utilities.GetEnv(utilities.GetEnvOptions{
		DefaultValue: fmt.Sprint(constants.TOKENS.DefaultRefreshTokenExpirationSeconds),
		EnvName:      constants.ENV_NAMES.RefreshTokenExpirationSeconds,
	})
	refreshTokenExpirationSeconds, convertError := strconv.Atoi(
		refreshTokenExpirationSecondsString,
	)
	if convertError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: convertError,
		})
	}

	queryError = tx.Create(&postgresql.UsedRefreshToken{
		ExpiresAt: issuedAtSeconds + int64(refreshTokenExpirationSeconds),
		Token:     refreshToken,
		UserID:    uint(userId),
	}).Error
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
	setError := redis.Client.Set(
		context.Context(),
		redis.CreateKey(
			constants.REDIS_PREFIXES.BlacklistedTokenPair,
			fmt.Sprintf("%d-%s", userId, refreshTokenPairId),
		),
		refreshTokenPairId,
		time.Duration(accessTokenExpirationSeconds)*time.Second,
	).Err()
	if setError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: setError,
		})
	}

	return utilities.Response(utilities.ResponseOptions{Context: context})
}
