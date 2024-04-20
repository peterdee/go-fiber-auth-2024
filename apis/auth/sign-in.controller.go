package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/julyskies/gohelpers"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/utilities"
)

func signInController(context fiber.Ctx) error {
	if len(context.Body()) == 0 {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	payload := new(SignInPayload)
	if bindingError := context.Bind().Body(payload); bindingError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: bindingError,
		})
	}
	if payload.Email == "" || payload.Password == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	email := strings.Trim(strings.ToLower(payload.Email), " ")

	var userRecord postgresql.User
	queryError := postgresql.
		Database.
		Where(&postgresql.User{Email: email}).
		First(&userRecord).
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

	var passwordRecord postgresql.Password
	queryError = postgresql.
		Database.
		Where(&postgresql.Password{UserID: userRecord.ID}).
		First(&passwordRecord).
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

	password := strings.Trim(payload.Password, " ")
	isValid, matchError := utilities.ComparePasswordAndHash(password, passwordRecord.Hash)
	if matchError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: matchError,
		})
	}
	if !isValid {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.Unauthorized,
			Status: fiber.StatusUnauthorized,
		})
	}

	var userSecretRecord postgresql.Password
	queryError = postgresql.
		Database.
		Where(&postgresql.UserSecret{UserID: userRecord.ID}).
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

	tokenPairId := gohelpers.RandomString(24)

	accessTokenSecret, secretError := utilities.CreateTokenSecret(
		userSecretRecord.Hash,
		passwordRecord.Hash,
		utilities.GetEnv(utilities.GetEnvOptions{
			DefaultValue: constants.TOKENS.DefaultAccessTokenCommonSecret,
			EnvName:      constants.ENV_NAMES.AccessTokenCommonSecret,
		}),
	)
	if secretError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: secretError,
		})
	}
	accessToken, tokenError := utilities.CreateToken(
		fmt.Sprint(userRecord.ID),
		accessTokenSecret,
		tokenPairId,
	)
	if tokenError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: tokenError,
		})
	}

	refreshTokenSecret, secretError := utilities.CreateTokenSecret(
		userSecretRecord.Hash,
		passwordRecord.Hash,
		utilities.GetEnv(utilities.GetEnvOptions{
			DefaultValue: constants.TOKENS.DefaultRefreshTokenCommonSecret,
			EnvName:      constants.ENV_NAMES.RefreshTokenCommonSecret,
		}),
	)
	if secretError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: secretError,
		})
	}
	refreshToken, tokenError := utilities.CreateToken(
		fmt.Sprint(userRecord.ID),
		refreshTokenSecret,
		tokenPairId,
	)
	if tokenError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: tokenError,
		})
	}

	return utilities.Response(utilities.ResponseOptions{
		Context: context,
		Data: fiber.Map{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}
