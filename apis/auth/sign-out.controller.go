package auth

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
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

	refreshToken := strings.Trim(strings.ToLower(payload.RefreshToken), " ")

	tx := postgresql.Database.Begin()
	if tx.Error != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: tx.Error,
		})
	}

	var refreshTokenRecord postgresql.UsedRefreshToken
	queryError := postgresql.
		Database.
		Where(&postgresql.UsedRefreshToken{Token: refreshToken}). // TODO: use UserID in query
		First(&refreshTokenRecord).
		Error
	if queryError != nil && queryError.Error() != "record not found" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}
	if refreshTokenRecord.ID != 0 {
		// TODO: complete sign out by changing the secret
	}

	return utilities.Response(utilities.ResponseOptions{Context: context})
}
