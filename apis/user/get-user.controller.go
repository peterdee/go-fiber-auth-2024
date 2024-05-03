package user

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/postgresql"
	"go-fiber-auth-2024/utilities"
)

func getUser(context fiber.Ctx) error {
	userId := context.Params("id", "")
	if userId == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	var user postgresql.User
	queryError := postgresql.Database.Where("id = ?", userId).First(&user).Error
	if queryError != nil {
		if queryError.Error() == "record not found" {
			return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
				Info:   constants.RESPONSE_INFO.InvalidUserId,
				Status: fiber.StatusBadRequest,
			})
		}
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: queryError,
		})
	}

	return utilities.Response(utilities.ResponseOptions{
		Context: context,
		Data: fiber.Map{
			"user": user,
		},
	})
}
