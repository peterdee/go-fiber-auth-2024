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

	tx := postgresql.Database.Begin()
	if tx.Error != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: tx.Error,
		})
	}

	email := strings.Trim(strings.ToLower(payload.Email), " ")
	var existingUserRecord postgresql.User
	queryError := tx.
		Where(&postgresql.User{Email: email}).
		First(&existingUserRecord).
		Error
	if queryError != nil && queryError.Error() != "record not found" {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}
	if existingUserRecord.ID != 0 {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.EmailAlreadyInUse,
			Status: fiber.StatusBadRequest,
		})
	}

	firstName := strings.Trim(payload.FirstName, " ")
	lastName := strings.Trim(payload.LastName, " ")
	newUser := postgresql.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
	result := tx.Create(&newUser)
	if result.Error != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: result.Error,
		})
	}

	password := strings.Trim(payload.Password, " ")
	hashed, hashError := utilities.CreateHash(password)
	if hashError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: hashError,
		})
	}
	result = tx.Create(&postgresql.Password{
		Hash:   hashed,
		UserID: newUser.ID,
	})
	if result.Error != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: result.Error,
		})
	}

	secret, secretError := utilities.CreateUserSecret(newUser.ID)
	if secretError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: secretError,
		})
	}
	result = tx.Create(&postgresql.UserSecret{
		Secret: secret,
		UserID: newUser.ID,
	})
	if result.Error != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: result.Error,
		})
	}

	commitError := tx.Commit().Error
	if commitError != nil {
		tx.Rollback()
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: commitError,
		})
	}

	return utilities.Response(utilities.ResponseOptions{Context: context})
}
