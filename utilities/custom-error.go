package utilities

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
)

type ApplicationError struct {
	Err       error
	Info      string
	Printable bool
	Status    int
}

func (err *ApplicationError) Error() string {
	return err.Err.Error()
}

type ApplicationErrorOptions struct {
	Err    error
	Info   string
	Status int
}

func NewApplicationError(options ApplicationErrorOptions) *ApplicationError {
	newCustomError := new(ApplicationError)
	newCustomError.Err = errors.New(constants.RESPONSE_INFO.InternalServerError)
	newCustomError.Info = constants.RESPONSE_INFO.InternalServerError
	newCustomError.Printable = false
	newCustomError.Status = fiber.StatusInternalServerError

	if options.Info != "" {
		newCustomError.Info = options.Info
		newCustomError.Err = errors.New(options.Info)
	}

	if options.Status != 0 {
		newCustomError.Status = options.Status
	}

	if options.Err != nil {
		newCustomError.Err = options.Err
		newCustomError.Printable = true
	}

	return newCustomError
}
