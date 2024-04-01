package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

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
	password := strings.Trim(payload.Password, " ")

	// TODO: use controlled transaction
	var user postgresql.User

	tx := postgresql.Database.Begin()
	t
	transactionError := postgresql.Database.Transaction(
		func(tx *gorm.DB) error {
			queryError := tx.
				Where(&postgresql.User{Email: email}).
				First(&user).
				Error
			if queryError != nil {
				return queryError
			}
			if user.ID != 0 {

			}
			return nil
		},
	)

	if transactionError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: transactionError,
		})
	}
	fmt.Println("user", user)

	return utilities.Response(utilities.ResponseOptions{Context: context})
}
