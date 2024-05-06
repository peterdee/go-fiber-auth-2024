package redis

import (
	"context"
	"fmt"
	"log"
	"time"

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

	for i := 0; i <= 5; i += 1 {
		ping := Client.Ping(context.Background())
		if ping.Err() != nil {
			if i == 5 {
				log.Fatal(ping.Err())
			}
			log.Printf("Could not connect to Redis, retrying in %d sec", i+1)
			time.Sleep(time.Second * time.Duration(i+1))
		} else {
			break
		}
	}

	log.Println(constants.ACTION_MESSAGES.RedisConnected)
}

func CreateKey(prefix string, userId string) string {
	return fmt.Sprintf("%s-%s", prefix, userId)
}
