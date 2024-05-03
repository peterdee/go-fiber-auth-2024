package auth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/utilities"
)

func refreshTokensController(context fiber.Ctx) error {
	if len(context.Body()) == 0 {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	payload := new(RefreshTokensPayload)
	if bindingError := context.Bind().Body(payload); bindingError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: bindingError,
		})
	}
	if payload.AccessToken == "" || payload.RefreshToken == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	accessToken := strings.Trim(payload.AccessToken, " ")
	refreshToken := strings.Trim(payload.RefreshToken, " ")

	accessTokenClaims, decodeError := utilities.DecodeToken(accessToken)
	if decodeError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: decodeError,
		})
	}
	refreshTokenClaims, decodeError := utilities.DecodeToken(refreshToken)
	if decodeError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: decodeError,
		})
	}

	accessTokenUserId := accessTokenClaims.Subject
	refreshTokenUserId := refreshTokenClaims.Subject
	if accessTokenUserId != refreshTokenUserId {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.Unauthorized,
			Status: fiber.StatusUnauthorized,
		})
	}

	tx := postgresql.Database.Begin()
	if tx.Error != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: tx.Error,
		})
	}

	var user postgresql.User
	queryError := tx.Where("id = ?", accessTokenUserId).First(&user).Error
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

	var usedRefreshToken postgresql.UsedRefreshToken
	queryError = tx.Where("token = ?", refreshToken).First(&usedRefreshToken).Error
	if queryError != nil && queryError.Error() != "record not found" {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}
	if usedRefreshToken.ID != 0 {
		internalError := signOutFromAllDevices(int(user.ID), tx, context)
		if internalError != nil {
			tx.Rollback()
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: internalError,
			})
		}
		commitError := tx.Commit().Error
		if commitError != nil {
			tx.Rollback()
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: commitError,
			})
		}
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.Unauthorized,
			Status: fiber.StatusUnauthorized,
		})
	}

	accessTokenPairId := accessTokenClaims.ID
	refreshTokenPairId := refreshTokenClaims.ID
	if accessTokenPairId != refreshTokenPairId {
		internalError := signOutFromAllDevices(int(user.ID), tx, context)
		if internalError != nil {
			tx.Rollback()
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: internalError,
			})
		}
		commitError := tx.Commit().Error
		if commitError != nil {
			tx.Rollback()
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Err: commitError,
			})
		}
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.Unauthorized,
			Status: fiber.StatusUnauthorized,
		})
	}

	isExpired, internalError := utilities.CheckTokenExpiration(
		refreshTokenClaims.Expires.Time().Unix(),
		utilities.TOKEN_TYPE_REFRESH,
	)
	if internalError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: internalError,
		})
	}
	if isExpired {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.Unauthorized,
			Status: fiber.StatusUnauthorized,
		})
	}

	fingerprint := utilities.Fingerprint(context)
	var userPasswordRecord postgresql.Password
	queryError = tx.
		Where(&postgresql.Password{UserID: user.ID}).
		First(&userPasswordRecord).
		Error
	if queryError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.Unauthorized,
			Status: fiber.StatusUnauthorized,
		})
	}
	var userSecretRecord postgresql.UserSecret
	queryError = tx.
		Where(&postgresql.UserSecret{UserID: user.ID}).
		First(&userSecretRecord).
		Error
	if queryError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.Unauthorized,
			Status: fiber.StatusUnauthorized,
		})
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
	isValid := utilities.VerifyToken(refreshToken, refreshTokenSecret)
	if !isValid {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.Unauthorized,
			Status: fiber.StatusUnauthorized,
		})
	}

	tokenExpirationString := utilities.GetEnv(utilities.GetEnvOptions{
		DefaultValue: fmt.Sprint(constants.TOKENS.DefaultRefreshTokenExpirationSeconds),
		EnvName:      constants.ENV_NAMES.RefreshTokenExpirationSeconds,
	})
	tokenExpiration, convertError := strconv.Atoi(tokenExpirationString)
	if convertError != nil {
		tokenExpiration = constants.TOKENS.DefaultRefreshTokenExpirationSeconds
	}
	queryError = tx.Create(&postgresql.UsedRefreshToken{
		ExpiresAt: refreshTokenClaims.Expires.Time().Unix() + int64(tokenExpiration),
		Token:     refreshToken,
		UserID:    user.ID,
	}).Error
	if queryError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}

	newAccessToken, newRefreshToken, internalError := createTokens(
		user.ID,
		tx,
		context,
	)
	if internalError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: internalError,
		})
	}

	return utilities.Response(utilities.ResponseOptions{
		Context: context,
		Data: fiber.Map{
			"accessToken":  newAccessToken,
			"refreshToken": newRefreshToken,
		},
	})
}
