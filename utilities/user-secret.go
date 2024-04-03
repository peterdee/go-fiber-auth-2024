package utilities

import (
	"fmt"

	"github.com/julyskies/gohelpers"
)

func CreateUserSecret(userId uint) (string, error) {
	return CreateHash(fmt.Sprintf(
		"%d:%d:%s",
		userId,
		gohelpers.MakeTimestamp(),
		gohelpers.RandomString(16),
	))
}
