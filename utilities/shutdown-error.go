package utilities

import (
	"log"

	"go-fiber-auth-2024/constants"
)

func ShutdownError(shutdownError error) {
	log.Printf("%s:\n%s", constants.ACTION_MESSAGES.ShutdownError, shutdownError.Error())
}
