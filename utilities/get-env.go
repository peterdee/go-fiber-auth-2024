package utilities

import (
	"log"
	"os"

	"go-fiber-auth-2024/constants"
)

func GetEnv(options GetEnvOptions) string {
	value := os.Getenv(options.EnvName)
	if value == "" {
		if options.IsRequired {
			log.Fatalf("%s: %s", constants.ACTION_MESSAGES.EnvLoadingError, options.EnvName)
		}
		value = options.DefaultValue
	}
	return value
}
