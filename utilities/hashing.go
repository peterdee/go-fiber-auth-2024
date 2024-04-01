package utilities

import (
	"errors"

	"github.com/alexedwards/argon2id"
)

const EMPTY_HASHED_STRING_ERROR string = "provided hashed string is empty"
const EMPTY_PLAINTEXT_STRING_ERROR string = "provided plaintext string is empty"

func ComparePasswordAndHash(plaintext, hashed string) (bool, error) {
	if hashed == "" {
		return false, errors.New(EMPTY_HASHED_STRING_ERROR)
	}
	if plaintext == "" {
		return false, errors.New(EMPTY_PLAINTEXT_STRING_ERROR)
	}

	isValid, compareError := argon2id.ComparePasswordAndHash(plaintext, hashed)
	if compareError != nil {
		return false, compareError
	}
	return isValid, nil
}

func CreateHash(plaintext string) (string, error) {
	if plaintext == "" {
		return "", errors.New(EMPTY_PLAINTEXT_STRING_ERROR)
	}

	hash, hashError := argon2id.CreateHash(plaintext, argon2id.DefaultParams)
	if hashError != nil {
		return "", hashError
	}
	return hash, nil
}
