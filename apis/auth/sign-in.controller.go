package auth

import (
	"strings"

	"github.com/gofiber/fiber/v3"

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

	return utilities.Response(utilities.ResponseOptions{Context: context})
}
