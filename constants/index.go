package constants

var ACTION_MESSAGES = ActionMessages{
	EnvLoadingError:       "Could not load ENV variable",
	LoadedEnvFile:         "Loaded .env file",
	PGConnected:           "PostgreSQL connected",
	PGConnectionClosed:    "PostgreSQL connection closed",
	RedisConnected:        "Redis connected",
	ShutdownError:         "Error while gracefully shutting down the server",
	ShutdownSuccess:       "Server gracefully stopped",
	UserIDAssertionFailed: "userId assertion failed",
}

const APP_NAME string = "GO FIBER AUTH 2024"

const DEFAULT_REDIS_HOST string = "localhost:6379"

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
	InvlaidToken:           "INVALID_TOKEN",
	MissingData:            "MISSING_DATA",
	MissingToken:           "MISSING_TOKEN",
	Ok:                     "OK",
	Unauthorized:           "UNAUTHORIZED",
}

var TOKENS = Tokens{
	DefaultAccessTokenCommonSecret:       "access-secret",
	DefaultAccessTokenExpirationSeconds:  3600,
	DefaultRefreshTokenCommonSecret:      "refresh-secret",
	DefaultRefreshTokenExpirationSeconds: 36000,
}
