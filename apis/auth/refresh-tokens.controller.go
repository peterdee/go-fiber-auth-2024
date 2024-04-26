package auth

import (
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

	// get user ids from AT & RT
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

	// check if user record exists
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

	// compare token pair ids
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

	// check refresh token expiration
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

	// TODO: verify refresh token, store it in the database, create a new token pair and return it

	return utilities.Response(utilities.ResponseOptions{Context: context})
}
