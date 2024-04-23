package constants

type ActionMessages struct {
	EnvLoadingError       string
	LoadedEnvFile         string
	PGConnected           string
	PGConnectionClosed    string
	RedisConnected        string
	ShutdownError         string
	ShutdownSuccess       string
	UserIDAssertionFailed string
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
	RedisHost                     string
	RedisPassword                 string
	RedisUsername                 string
	RefreshTokenCommonSecret      string
	RefreshTokenExpirationSeconds string
}

type EnvSources struct {
	Env  string
	File string
}

type LocalsKeys struct {
	RequestReceivedTimestamp string
	TokenPairId              string
	UserId                   string
}

type RedisPrefixes struct {
	PasswordHash string
	SecretHash   string
}

type ResponseInfo struct {
	AccessTokenExpired  string
	EmailAlreadyInUse   string
	InternalServerError string
	InvlaidToken        string
	MissingData         string
	MissingToken        string
	Ok                  string
	Unauthorized        string
}

type Tokens struct {
	DefaultAccessTokenCommonSecret       string
	DefaultAccessTokenExpirationSeconds  int
	DefaultRefreshTokenCommonSecret      string
	DefaultRefreshTokenExpirationSeconds int
}
