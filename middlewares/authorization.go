package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/utilities"
)

func Authorization(context fiber.Ctx) error {
	accessTokenRaw := context.Get("authorization")
	if accessTokenRaw == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingToken,
			Status: fiber.StatusUnauthorized,
		})
	}
	accessToken := strings.Trim(accessTokenRaw, " ")
	if accessToken == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.MissingToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	claims, decodeError := utilities.DecodeToken(accessToken)
	if decodeError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.InvlaidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	issuedAt := claims.Issued
	tokenPairId := claims.ID
	userIdRaw := claims.Subject

	if issuedAt == nil || tokenPairId == "" || userIdRaw == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.InvlaidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	// TODO: check expiration, validate token & proceed

	context.Locals(
		constants.LOCALS_KEYS.UserId,
		0,
	)

	return context.Next()
}
