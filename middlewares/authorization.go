package middlewares

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/julyskies/gohelpers"

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

	issuedAt := claims.Issued.Time().UnixNano() / int64(time.Millisecond)
	tokenPairId := claims.ID
	userIdRaw := claims.Subject

	if issuedAt == 0 || tokenPairId == "" || userIdRaw == "" {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Info:   constants.RESPONSE_INFO.InvlaidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	accessTokenExpirationString := utilities.GetEnv(utilities.GetEnvOptions{
		DefaultValue: fmt.Sprint(constants.TOKENS.DefaultAccessTokenExpirationSeconds),
		EnvName:      constants.ENV_NAMES.AccessTokenExpirationSeconds,
	})
	accessTokenExpiration, convertError := strconv.Atoi(accessTokenExpirationString)
	if convertError != nil {
		return utilities.NewApplicationError(utilities.ApplicationErrorOptions{
			Err: convertError,
		})
	}
	if issuedAt+int64(accessTokenExpiration) < gohelpers.MakeTimestamp() {
		// TODO: expiration error
	}

	// TODO: validate token & proceed

	context.Locals(
		constants.LOCALS_KEYS.UserId,
		0,
	)

	return context.Next()
}
