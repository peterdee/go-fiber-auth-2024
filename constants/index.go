package constants

var ACTION_MESSAGES = ActionMessages{
	EnvLoadingError:     "Could not load ENV variable",
	InvalidTokenType:    "invalid token type",
	LoadedEnvFile:       "Loaded .env file",
	PGConnected:         "PostgreSQL connected",
	PGConnectionClosed:  "PostgreSQL connection closed",
	RedisConnected:      "Redis connected",
	ShutdownError:       "Error while gracefully shutting down the server",
	ShutdownSuccess:     "Server gracefully stopped",
	TypeAssertionFailed: "type assertion failed",
}

const APP_NAME string = "GO FIBER AUTH 2024"

const DEFAULT_REDIS_HOST string = "localhost"

const DEFAULT_REDIS_PORT uint = 6379

var ENV_NAMES = EnvNames{
	AccessTokenCommonSecret:       "ACCESS_TOKEN_COMMON_SECRET",
	AccessTokenExpirationSeconds:  "ACCESS_TOKEN_EXPIRATION_SECONDS",
	EnvSource:                     "ENV_SOURCE",
	PGDatabase:                    "PG_DATABASE",
	PGHost:                        "PG_HOST",
	PGPassword:                    "PG_PASSWORD",
	PGPort:                        "PG_PORT",
	PGUsername:                    "PG_USERNAME",
	Port:                          "PORT",
	RedisHost:                     "REDIS_HOST",
	RedisPassword:                 "REDIS_PASSWORD",
	RedisPort:                     "REDIS_PORT",
	RedisUsername:                 "REDIS_USERNAME",
	RefreshTokenCommonSecret:      "REFRESH_TOKEN_COMMON_SECRET",
	RefreshTokenExpirationSeconds: "REFRESH_TOKEN_EXPIRATION_SECONDS",
}

var ENV_SOURCES = EnvSources{
	Env:  "env",
	File: "file",
}

var LOCALS_KEYS = LocalsKeys{
	RequestReceivedTimestamp: "requestReceivedTimestamp",
	TokenPairId:              "tokenPairId",
	UserId:                   "userId",
}

const PORT string = "2024"

var REDIS_PREFIXES = RedisPrefixes{
	BlacklistedTokenPair: "bltpair",
	PasswordHash:         "phash",
	SecretHash:           "shash",
}

var RESPONSE_INFO = ResponseInfo{
	AccessTokenBlacklisted: "ACCESS_TOKEN_BLACKLISTED",
	AccessTokenExpired:     "ACCESS_TOKEN_EXPIRED",
	EmailAlreadyInUse:      "EMAIL_ALREADY_IN_USE",
	InternalServerError:    "INTERNAL_SERVER_ERROR",
	InvalidToken:           "INVALID_TOKEN",
	InvalidUserId:          "INVALID_USER_ID",
	MissingData:            "MISSING_DATA",
	MissingToken:           "MISSING_TOKEN",
	Ok:                     "OK",
	OldPasswordIsInvalid:   "OLD_PASSWORD_IS_INVALID",
	Unauthorized:           "UNAUTHORIZED",
}

const TOKEN_PAIR_ID_LENGTH int = 24

var TOKENS = Tokens{
	DefaultAccessTokenCommonSecret:       "access-secret",
	DefaultAccessTokenExpirationSeconds:  3600,
	DefaultRefreshTokenCommonSecret:      "refresh-secret",
	DefaultRefreshTokenExpirationSeconds: 36000,
}
