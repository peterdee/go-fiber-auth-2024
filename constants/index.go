package constants

var ACTION_MESSAGES = ActionMessages{
	LoadedEnvFile:      "Loaded .env file",
	PGConnected:        "PostgreSQL connected",
	PGCredentialsError: "Could not load PostgreSQL credentials",
	ShutdownError:      "Error while gracefully shutting down the server",
	ShutdownSuccess:    "Server gracefully stopped",
}

const APP_NAME string = "GO FIBER AUTH 2024"

var ENV_NAMES = EnvNames{
	EnvSource:  "ENV_SOURCE",
	Port:       "PORT",
	PGDatabase: "PG_DATABASE",
	PGHost:     "PG_HOST",
	PGPassword: "PG_PASSWORD",
	PGPort:     "PG_PORT",
	PGUsername: "PG_USERNAME",
}

var ENV_SOURCES = EnvSources{
	Env:  "env",
	File: "file",
}

var LOCALS_KEYS = LocalsKeys{
	RequestReceivedTimestamp: "requestReceivedTimestamp",
	UserId:                   "userId",
}

const PORT string = "2024"

var RESPONSE_INFO = ResponseInfo{
	InternalServerError: "INTERNAL_SERVER_ERROR",
	Ok:                  "OK",
}
