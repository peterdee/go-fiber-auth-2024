package utilities

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/julyskies/gohelpers"
	"github.com/pascaldekloe/jwt"

	"go-fiber-auth-2024/constants"
)

type extraHeaders struct {
	Type string `json:"typ"`
}

const TOKEN_TYPE_ACCESS string = "access"
const TOKEN_TYPE_REFRESH string = "refresh"

func CheckTokenExpiration(issuedAtSeconds int64, tokenType string) (bool, error) {
	if tokenType != TOKEN_TYPE_ACCESS && tokenType != TOKEN_TYPE_REFRESH {
		return false, errors.New(constants.ACTION_MESSAGES.InvalidTokenType)
	}

	var tokenExpirationString string
	if tokenType == TOKEN_TYPE_ACCESS {
		tokenExpirationString = GetEnv(GetEnvOptions{
			DefaultValue: fmt.Sprint(constants.TOKENS.DefaultAccessTokenExpirationSeconds),
			EnvName:      constants.ENV_NAMES.AccessTokenExpirationSeconds,
		})
	}
	if tokenType == TOKEN_TYPE_REFRESH {
		tokenExpirationString = GetEnv(GetEnvOptions{
			DefaultValue: fmt.Sprint(constants.TOKENS.DefaultRefreshTokenExpirationSeconds),
			EnvName:      constants.ENV_NAMES.RefreshTokenExpirationSeconds,
		})
	}

	tokenExpiration, convertError := strconv.Atoi(tokenExpirationString)
	if convertError != nil {
		if tokenType == TOKEN_TYPE_REFRESH {
			tokenExpiration = constants.TOKENS.DefaultRefreshTokenExpirationSeconds
		} else {
			tokenExpiration = constants.TOKENS.DefaultAccessTokenExpirationSeconds
		}
	}
	if issuedAtSeconds+int64(tokenExpiration) < gohelpers.MakeTimestampSeconds() {
		return true, nil
	}
	return false, nil
}

func CreateTokenSecret(
	userSecretHash,
	userPasswordHash,
	commonSecret,
	fingerprint string,
) string {
	hashed := md5.Sum([]byte(fmt.Sprintf(
		"%s:%s:%s:%s",
		userSecretHash,
		userPasswordHash,
		commonSecret,
		fingerprint,
	)))
	return hex.EncodeToString(hashed[:])
}

func CreateToken(userId string, tokenSecret string, pairId string) (string, error) {
	var claims jwt.Claims

	claims.ID = pairId
	claims.Issued = jwt.NewNumericTime(time.Now().Round(time.Second))
	claims.Subject = userId

	var extra extraHeaders
	extra.Type = "JWT"
	extraJSON, jsonError := json.Marshal(extra)
	if jsonError != nil {
		return "", jsonError
	}

	token, signError := claims.HMACSign(jwt.HS256, []byte(tokenSecret), extraJSON)
	if signError != nil {
		return "", signError
	}
	return string(token), nil
}

func DecodeToken(token string) (*jwt.Claims, error) {
	invalidToken := errors.New("invalid token")

	claims, parseError := jwt.ParseWithoutCheck([]byte(token))
	if parseError != nil {
		return nil, parseError
	}
	if claims.Issued == nil {
		return nil, invalidToken
	}
	if claims.Subject == "" {
		return nil, invalidToken
	}

	return claims, nil
}

func VerifyToken(token, tokenSecret string) bool {
	_, verificationError := jwt.HMACCheck([]byte(token), []byte(tokenSecret))
	return verificationError == nil
}
