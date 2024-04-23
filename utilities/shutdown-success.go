package utilities

import (
	"log"

	"go-fiber-auth-2024/constants"
)

func ShutdownSuccess() {
	log.Println(constants.ACTION_MESSAGES.ShutdownSuccess)
}
