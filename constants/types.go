package constants

type ActionMessages struct {
	LoadedEnvFile      string
	PGConnected        string
	PGCredentialsError string
	ShutdownError      string
	ShutdownSuccess    string
}
type EnvNames struct {
	EnvSource  string
	Port       string
	PGDatabase string
	PGHost     string
	PGPassword string
	PGPort     string
	PGUsername string
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
	InternalServerError string
	Ok                  string
}
