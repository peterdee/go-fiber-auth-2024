package auth

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/utilities"
)

func signUpController(context fiber.Ctx) error {
	if len(context.Body()) == 0 {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	payload := new(SignUpPayload)
	if bindingError := context.Bind().Body(payload); bindingError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: bindingError,
		})
	}
	if payload.Email == "" || payload.FirstName == "" ||
		payload.LastName == "" || payload.Password == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	// TODO: use transaction
	email := strings.Trim(strings.ToLower(payload.Email), " ")
	var existingUserRecord postgresql.User
	queryError := postgresql.
		Database.
		Where(&postgresql.User{Email: email}).
		First(&existingUserRecord).
		Error
	if queryError != nil && queryError.Error() != "record not found" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}
	if existingUserRecord.ID != 0 {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.EmailAlreadyInUse,
			Status: fiber.StatusBadRequest,
		})
	}

	return utilities.Response(utilities.ResponseOptions{Context: context})
}
