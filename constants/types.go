package constants

type ActionMessages struct {
	LoadedEnvFile      string
	PGConnected        string
	PGCredentialsError string
	ShutdownError      string
	ShutdownSuccess    string
}
type EnvNames struct {
	AccessTokenCommonSecret       string
	AccessTokenExpirationSeconds  string
	EnvSource                     string
	PGDatabase                    string
	PGHost                        string
	PGPassword                    string
	PGPort                        string
	PGUsername                    string
	Port                          string
	RefreshTokenCommonSecret      string
	RefreshTokenExpirationSeconds string
}

type EnvSources struct {
	Env  string
	File string
}

type LocalsKeys struct {
	RequestReceivedTimestamp string
	UserId                   string
}

type ResponseInfo struct {
	EmailAlreadyInUse   string
	InternalServerError string
	MissingData         string
	Ok                  string
	Unauthorized        string
}
