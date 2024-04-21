package constants

type ActionMessages struct {
	EnvLoadingError string
	LoadedEnvFile   string
	PGConnected     string
	ShutdownError   string
	ShutdownSuccess string
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

type Tokens struct {
	DefaultAccessTokenCommonSecret       string
	DefaultAccessTokenExpirationSeconds  int
	DefaultRefreshTokenCommonSecret      string
	DefaultRefreshTokenExpirationSeconds int
}
