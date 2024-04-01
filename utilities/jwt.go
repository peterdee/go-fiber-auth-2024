package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pascaldekloe/jwt"
)

type extraHeaders struct {
	Type string `json:"typ"`
}

type TokenClaims struct {
	Issued  jwt.NumericTime `json:"iat"`
	Subject uint            `json:"sub"`
}

func CreateSecret(userSecretHash, passwordHash string) (string, error) {
	hashed, hashError := CreateHash(fmt.Sprintf("%s:%s", userSecretHash, passwordHash))
	if hashError != nil {
		return "", hashError
	}
	return hashed, nil
}

func CreateToken(userId string, tokenSecret string) (string, error) {
	var claims jwt.Claims

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
