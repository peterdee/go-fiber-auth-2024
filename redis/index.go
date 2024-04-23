package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	"go-fiber-auth-2024/constants"
	"go-fiber-auth-2024/utilities"
)

var Client *redis.Client

func CreateDatabaseConnection() {
	host := utilities.GetEnv(utilities.GetEnvOptions{
		DefaultValue: constants.DEFAULT_REDIS_HOST,
		EnvName:      constants.ENV_NAMES.RedisHost,
	})
	password := utilities.GetEnv(utilities.GetEnvOptions{
		EnvName: constants.ENV_NAMES.RedisPassword,
	})
	username := utilities.GetEnv(utilities.GetEnvOptions{
		EnvName: constants.ENV_NAMES.RedisUsername,
	})

	Client = redis.NewClient(&redis.Options{
		Addr:     host,
		DB:       0,
		Password: password,
		Username: username,
	})

	ping := Client.Ping(context.Background())
	if ping.Err() != nil {
		log.Fatal(ping.Err())
	}

	log.Println(constants.ACTION_MESSAGES.RedisConnected)
}
